package main

import (
	"log"
	"os"
	"time"

	"github.com/woningfinder/woningfinder/internal/corporation/scheduler"

	"github.com/woningfinder/woningfinder/internal/services/corporation"

	"github.com/woningfinder/woningfinder/internal/bootstrap"
	"github.com/woningfinder/woningfinder/pkg/config"
	"github.com/woningfinder/woningfinder/pkg/logging"

	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
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

	dbClient := bootstrap.CreateDBClient(logger)
	redisClient := bootstrap.CreateRedisClient(logger)
	mapboxClient := bootstrap.CreateMapboxClient()

	clientProvider := bootstrap.CreateClientProvider(logger, mapboxClient)
	corporationService := corporation.NewService(logger, dbClient, redisClient)

	// get time location
	nl, err := time.LoadLocation("Europe/Amsterdam")
	if err != nil {
		logger.Sugar().Fatal(err)
	}

	// instantiate cron
	c := cron.New(cron.WithLocation(nl), cron.WithSeconds(), cron.WithLogger(cron.VerbosePrintfLogger(log.New(os.Stdout, "cron: ", log.LstdFlags))))

	// populate crons
	for _, corp := range *clientProvider.List() {
		corp := corp // https://github.com/golang/go/wiki/CommonMistakes#using-reference-to-loop-iterator-variable

		// get corporation client
		client, err := clientProvider.Get(corp)
		if err != nil {
			logger.Sugar().Error(err)
			continue
		}

		// schedule corporation fetching
		schedule := scheduler.CorporationScheduler(corp)
		for _, s := range schedule {
			c.Schedule(s, cron.FuncJob(func() {
				if err := corporationService.PublishOffers(client, corp); err != nil {
					logger.Sugar().Error(err)
				}
			}))
		}
	}

	// start cron scheduler
	c.Run()
}
