package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/julienrbrt/woningfinder/cmd/woningfinder/dbmigrations"
	"github.com/julienrbrt/woningfinder/cmd/woningfinder/job"
	"github.com/julienrbrt/woningfinder/internal/auth"
	"github.com/julienrbrt/woningfinder/internal/bootstrap"
	bootstrapCorporation "github.com/julienrbrt/woningfinder/internal/bootstrap/corporation"
	"github.com/julienrbrt/woningfinder/internal/corporation"
	"github.com/julienrbrt/woningfinder/internal/corporation/city"
	"github.com/julienrbrt/woningfinder/internal/corporation/connector"
	"github.com/julienrbrt/woningfinder/internal/customer/matcher"
	"github.com/julienrbrt/woningfinder/internal/database"
	"github.com/julienrbrt/woningfinder/internal/handler"
	corporationService "github.com/julienrbrt/woningfinder/internal/services/corporation"
	emailService "github.com/julienrbrt/woningfinder/internal/services/email"
	matcherService "github.com/julienrbrt/woningfinder/internal/services/matcher"
	userService "github.com/julienrbrt/woningfinder/internal/services/user"
	"github.com/julienrbrt/woningfinder/pkg/config"
	"github.com/julienrbrt/woningfinder/pkg/logging"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

func init() {
	// load env variables
	if err := godotenv.Load(); err != nil {
		_ = config.MustGetString("DATABASE_URL")
	}
}

func main() {
	logger := logging.NewZapLogger(config.GetBoolOrDefault("APP_DEBUG", false), config.GetStringOrDefault("SENTRY_DSN", ""))
	jwtAuth := auth.CreateJWTAuthenticationToken(config.MustGetString("JWT_SECRET"))

	dbClient := bootstrap.CreateDBClient(logger)
	mapboxClient := bootstrap.CreateMapboxClient(logger, dbClient)
	imgClient := bootstrap.CreateImgDownloader(logger)
	emailClient := bootstrap.CreateEmailClient()

	connectorProvider := bootstrapCorporation.CreateConnectorProvider(logger, mapboxClient)
	cityTable := connectorProvider.GetCities()
	citySuggester := city.NewSuggester(cityTable)
	corporationService := corporationService.NewService(logger, dbClient, citySuggester)
	userService := userService.NewService(logger, dbClient, config.MustGetString("AES_SECRET"), connectorProvider, corporationService)
	emailService := emailService.NewService(logger, emailClient, jwtAuth, imgClient)
	matcherService := matcherService.NewService(logger, dbClient, userService, emailService, corporationService, imgClient, matcher.NewMatcher(citySuggester), connectorProvider)

	// set location to the netherlands
	nl, err := time.LoadLocation("Europe/Amsterdam")
	if err != nil {
		logger.Fatal("error when loading location", zap.Error(err))
	}

	// run db migrations
	if err := dbmigrations.RunMigrations(logger, dbClient, len(os.Args) > 1 && os.Args[1] == "init"); err != nil {
		logger.Fatal("error when running migrations", zap.Error(err))
	}

	// update city location (when a city has been added automatically and then updated in code)
	if err := updateCityLocation(logger, dbClient, cityTable, connectorProvider, corporationService); err != nil {
		logger.Fatal("error when updating city location", zap.Error(err))
	}

	// housing matcher
	go func(ch chan corporation.Offers) {
		for offers := range ch {
			logger.Info("received offers", zap.String("corporation", offers.CorporationName), zap.Int("count", len(offers.Offer)))

			if err := matcherService.Match(context.Background(), offers); err != nil {
				logger.Error("error while maching offers", zap.String("corporation", offers.CorporationName), zap.Error(err))
			}
		}
	}(job.OffersChan)

	// instantiate job and cron
	job := job.NewJobs(logger, dbClient, userService, matcherService, emailService)
	c := cron.New(cron.WithLocation(nl), cron.WithSeconds(), cron.WithLogger(cron.VerbosePrintfLogger(log.New(os.Stdout, "cron: ", log.LstdFlags))))

	// populate crons
	job.CleanupUnconfirmedCustomer(c)
	job.HousingFinder(c, connectorProvider)
	job.SendWeeklyUpdate(c)
	job.SendCorporationCredentialsMissingReminder(c)

	// start cron scheduler in a new go routine
	c.Start()

	// instantiate http server
	router := handler.NewHandler(logger, jwtAuth, corporationService, userService, emailService, imgClient)
	if err := http.ListenAndServe(":"+config.MustGetString("APP_PORT"), router); err != nil {
		logger.Fatal("failed to start server", zap.Error(err))
	}
}

// updateCityLocation updates the city location in the database.
// it is used because cities are defined in the code as well.
func updateCityLocation(
	logger *logging.Logger,
	dbClient database.DBClient,
	cityTable map[string]city.City,
	connectorProvider connector.ConnectorProvider,
	corporationService corporationService.Service,
) error {
	var cities []city.City
	if err := dbClient.Conn().Model(&cities).Where("latitude IS NULL OR longitude IS NULL").Select(); err != nil {
		return fmt.Errorf("error while getting cities without location: %w", err)
	}

	for _, c := range cities {
		city, ok := cityTable[c.Name]
		if !ok {
			logger.Warn("failed finding city", zap.String("city", c.Name))
			continue
		}

		logger.Info("updating city", zap.String("city", city.Name))
		if _, err := dbClient.Conn().Model(&city).OnConflict("(name) DO UPDATE").Insert(); err != nil {
			return fmt.Errorf("failed updating city %s: %w", city.Name, err)
		}
	}

	// update corporations cities
	for _, corp := range connectorProvider.GetCorporations() {
		if err := corporationService.LinkCities(corp.Cities, true, corp); err != nil {
			return fmt.Errorf("failed updating corporation %s: %w", corp.Name, err)
		}
	}

	return nil
}
