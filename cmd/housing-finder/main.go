package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/woningfinder/woningfinder/internal/bootstrap"
	"github.com/woningfinder/woningfinder/internal/corporation"
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
	logger := logging.NewZapLoggerWithSentry(config.MustGetString("SENTRY_DSN"))

	// connect to databases
	err := bootstrap.InitDB()
	if err != nil {
		logger.Sugar().Fatal(err)
	}

	err = bootstrap.InitRedis()
	if err != nil {
		logger.Sugar().Fatal(err)
	}

	mapboxClient := bootstrap.CreateMapboxGeocodingClient()
	clientProvider := bootstrap.CreateClientProvider(logger, mapboxClient)
	corporationService := corporation.NewService(logger, bootstrap.DB, bootstrap.RDB)
	// get time location
	nl, err := time.LoadLocation("Europe/Amsterdam")
	if err != nil {
		logger.Sugar().Fatal(err)
	}

	// instantiate cron
	c := cron.New(cron.WithLocation(nl), cron.WithSeconds(), cron.WithLogger(cron.VerbosePrintfLogger(log.New(os.Stdout, "cron: ", log.LstdFlags))))

	// populate crons
	for _, corp := range *clientProvider.List() {
		// get corporation client
		corporationClient, err := clientProvider.Get(corp)
		if err != nil {
			logger.Sugar().Error(err)
			continue
		}

		// specify cron time or fallback to default
		// sets a second retry cron at 5 seconds interval
		var spec string
		for _, second := range []int{0, 5} {
			if corp.SelectionTime == (time.Time{}) {
				// the default is running at 17:15 is not specified
				spec = buildSpec(17, 15, second)
			} else {
				spec = buildSpec(corp.SelectionTime.Hour(), corp.SelectionTime.Minute(), second)
			}

			_, err = c.AddFunc(spec, func() {
				if err := corporationService.PublishOffers(corporationClient, corp); err != nil {
					logger.Sugar().Error(err)
				}
			})
			if err != nil {
				logger.Sugar().Error(err)
			}
		}

		// add one more time at midnight
		_, err = c.AddFunc("@midnight", func() {
			if err := corporationService.PublishOffers(corporationClient, corp); err != nil {
				logger.Sugar().Error(err)
			}
		})
		if err != nil {
			logger.Sugar().Error(err)
		}
	}

	// start cron scheduler
	c.Run()
}

func buildSpec(hour, minute, second int) string {
	return fmt.Sprintf("%d %d %d * * *", second, minute, hour)
}
