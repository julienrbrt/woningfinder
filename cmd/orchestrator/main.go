package main

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/julienrbrt/woningfinder/cmd/orchestrator/job"
	"github.com/julienrbrt/woningfinder/internal/auth"
	"github.com/julienrbrt/woningfinder/internal/bootstrap"
	bootstrapCorporation "github.com/julienrbrt/woningfinder/internal/bootstrap/corporation"
	"github.com/julienrbrt/woningfinder/internal/corporation/city"
	"github.com/julienrbrt/woningfinder/internal/customer/matcher"
	corporationService "github.com/julienrbrt/woningfinder/internal/services/corporation"
	emailService "github.com/julienrbrt/woningfinder/internal/services/email"
	matcherService "github.com/julienrbrt/woningfinder/internal/services/matcher"
	userService "github.com/julienrbrt/woningfinder/internal/services/user"
	"github.com/julienrbrt/woningfinder/pkg/config"
	"github.com/julienrbrt/woningfinder/pkg/logging"
	"github.com/robfig/cron/v3"
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
	logger := logging.NewZapLogger(config.GetBoolOrDefault("APP_DEBUG", false), config.GetStringOrDefault("SENTRY_DSN", ""))
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
