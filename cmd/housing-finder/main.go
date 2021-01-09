package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/woningfinder/woningfinder/internal/bootstrap"
	"github.com/woningfinder/woningfinder/internal/config"
	"github.com/woningfinder/woningfinder/internal/corporation"
	"github.com/woningfinder/woningfinder/internal/logging"

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
	logger := logging.NewZapLoggerWithSentry(config.MustGetString("SENTRY_DSN"))

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
		// https://github.com/golang/go/wiki/CommonMistakes#using-reference-to-loop-iterator-variable
		corp := corp

		// get corporation client
		client, err := clientProvider.Get(corp)
		if err != nil {
			logger.Sugar().Error(err)
			continue
		}

		// specify cron time or fallback to default
		// always check at midnight
		if _, err = c.AddFunc("@midnight", func() {
			if err := corporationService.PublishOffers(client, corp); err != nil {
				logger.Sugar().Error(err)
			}
		}); err != nil {
			logger.Sugar().Error(err)
		}

		// check at 0, 10 and 30 seconds after the publishing time
		var spec string
		for _, second := range []int{0, 10, 25, 50} {
			if corp.SelectionTime != (time.Time{}) {
				spec = buildSpec(corp.SelectionTime.Hour(), corp.SelectionTime.Minute(), second)
			} else {
				// the default is running at 17:00 if not specified
				spec = buildSpec(17, 00, second)
			}

			if _, err = c.AddFunc(spec, func() {
				if err := corporationService.PublishOffers(client, corp); err != nil {
					logger.Sugar().Error(err)
				}
			}); err != nil {
				logger.Sugar().Error(err)
			}
		}
	}

	// start cron scheduler
	c.Run()
}

func buildSpec(hour, minute, second int) string {
	return fmt.Sprintf("%d %d %d * * *", second, minute, hour)
}
