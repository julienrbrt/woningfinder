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
	corporationService "github.com/woningfinder/woningfinder/internal/services/corporation"
	matcherService "github.com/woningfinder/woningfinder/internal/services/matcher"
	notificationsService "github.com/woningfinder/woningfinder/internal/services/notifications"
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

	// register m2m models with go-pg
	bootstrap.RegisterModel()
}

func main() {
	logger := logging.NewZapLogger(config.GetBoolOrDefault("APP_DEBUG", false), config.MustGetString("SENTRY_DSN"))
	jwtAuth := auth.CreateJWTAuthenticationToken(config.MustGetString("JWT_SECRET"))

	dbClient := bootstrap.CreateDBClient(logger)
	redisClient := bootstrap.CreateRedisClient(logger)
	mapboxClient := bootstrap.CreateMapboxClient()
	emailClient := bootstrap.CreateEmailClient()

	clientProvider := bootstrap.CreateClientProvider(logger, mapboxClient)
	corporationService := corporationService.NewService(logger, dbClient)
	userService := userService.NewService(logger, dbClient, redisClient, config.MustGetString("AES_SECRET"), clientProvider, corporationService)
	matcherService := matcherService.NewService(logger, redisClient, userService, corporationService, clientProvider)
	notificationsService := notificationsService.NewService(logger, emailClient, jwtAuth)

	// set location to the netherlands
	nl, err := time.LoadLocation("Europe/Amsterdam")
	if err != nil {
		logger.Sugar().Fatal(err)
	}

	// instantiate cron
	c := cron.New(cron.WithLocation(nl), cron.WithSeconds(), cron.WithLogger(cron.VerbosePrintfLogger(log.New(os.Stdout, "cron: ", log.LstdFlags))))

	// populate crons
	job.HousingFinder(logger, c, clientProvider, matcherService)
	job.BatchWeeklyUpdate(logger, c, userService, notificationsService)

	// start cron scheduler
	c.Run()
}