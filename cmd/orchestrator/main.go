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
	"github.com/woningfinder/woningfinder/internal/customer/matcher"
	corporationService "github.com/woningfinder/woningfinder/internal/services/corporation"
	matcherService "github.com/woningfinder/woningfinder/internal/services/matcher"
	notificationService "github.com/woningfinder/woningfinder/internal/services/notification"
	userService "github.com/woningfinder/woningfinder/internal/services/user"
	"github.com/woningfinder/woningfinder/pkg/config"
	"github.com/woningfinder/woningfinder/pkg/logging"
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
	mapboxClient := bootstrap.CreateMapboxClient()
	emailClient := bootstrap.CreateEmailClient()

	clientProvider := bootstrapCorporation.CreateClientProvider(logger, mapboxClient)
	corporationService := corporationService.NewService(logger, dbClient)
	userService := userService.NewService(logger, dbClient, redisClient, config.MustGetString("AES_SECRET"), clientProvider, corporationService)
	notificationService := notificationService.NewService(logger, emailClient, jwtAuth)
	matcherService := matcherService.NewService(logger, redisClient, userService, notificationService, corporationService, matcher.NewMatcher(), clientProvider)

	// set location to the netherlands
	nl, err := time.LoadLocation("Europe/Amsterdam")
	if err != nil {
		logger.Sugar().Fatal(err)
	}

	// instantiate job and cron
	job := job.NewJobs(logger, dbClient, redisClient, userService, matcherService, notificationService)
	c := cron.New(cron.WithLocation(nl), cron.WithSeconds(), cron.WithLogger(cron.VerbosePrintfLogger(log.New(os.Stdout, "cron: ", log.LstdFlags))))

	// populate crons
	job.CustomerAutoDelete(c)
	job.HousingFinder(c, clientProvider)
	job.SendCustomerPaymentReminder(c)
	job.SendWeeklyUpdate(c)

	// start cron scheduler
	c.Run()
}
