package main

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
	"github.com/woningfinder/woningfinder/cmd/orchestrator/job"
	"github.com/woningfinder/woningfinder/internal/auth"
	"github.com/woningfinder/woningfinder/internal/bootstrap"
	bootstrapCorporation "github.com/woningfinder/woningfinder/internal/bootstrap/corporation"
	"github.com/woningfinder/woningfinder/internal/corporation/city"
	"github.com/woningfinder/woningfinder/internal/customer/matcher"
	corporationService "github.com/woningfinder/woningfinder/internal/services/corporation"
	emailService "github.com/woningfinder/woningfinder/internal/services/email"
	matcherService "github.com/woningfinder/woningfinder/internal/services/matcher"
	userService "github.com/woningfinder/woningfinder/internal/services/user"
	"github.com/woningfinder/woningfinder/pkg/config"
	"github.com/woningfinder/woningfinder/pkg/logging"
	"go.uber.org/zap"
)

// init is invoked before main()
func init() {
	// loads values from .env into the system
	// fallback to system env if unexisting
	// if not defined on system, panics
	if err := godotenv.Load("../../.env"); err != nil {
		_ = config.MustGetString("APP_NAME")
	}
}

func main() {
	logger := logging.NewZapLogger(config.GetBoolOrDefault("APP_DEBUG", false), config.MustGetString("SENTRY_DSN"))
	jwtAuth := auth.CreateJWTAuthenticationToken(config.MustGetString("JWT_SECRET"))

	dbClient := bootstrap.CreateDBClient(logger)
	redisClient := bootstrap.CreateRedisClient(logger)
	mapboxClient := bootstrap.CreateMapboxClient(logger, redisClient)
	spacesClient := bootstrap.CreateDOSpacesClient(logger)
	emailClient := bootstrap.CreateEmailClient()

	connectorProvider := bootstrapCorporation.CreateConnectorProvider(logger, mapboxClient)
	citySuggester := city.NewSuggester(connectorProvider.GetCities())
	corporationService := corporationService.NewService(logger, dbClient, citySuggester)
	userService := userService.NewService(logger, dbClient, config.MustGetString("AES_SECRET"), connectorProvider, corporationService)
	emailService := emailService.NewService(logger, emailClient, jwtAuth)
	matcherService := matcherService.NewService(logger, redisClient, userService, emailService, corporationService, spacesClient, matcher.NewMatcher(citySuggester), connectorProvider)

	// set location to the netherlands
	nl, err := time.LoadLocation("Europe/Amsterdam")
	if err != nil {
		logger.Fatal("error when loading location", zap.Error(err))
	}

	// instantiate job and cron
	job := job.NewJobs(logger, dbClient, redisClient, userService, matcherService, emailService)
	c := cron.New(cron.WithLocation(nl), cron.WithSeconds(), cron.WithLogger(cron.VerbosePrintfLogger(log.New(os.Stdout, "cron: ", log.LstdFlags))))

	// populate crons
	job.CleanupUnconfirmedCustomer(c)
	job.HousingFinder(c, connectorProvider)
	job.SendWeeklyUpdate(c)
	job.SendCorporationCredentialsMissingReminder(c)

	// start cron scheduler
	c.Run()
}
