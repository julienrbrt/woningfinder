package main

import (
	"log"
	"os"
	"sync"
	"time"

	"github.com/woningfinder/woningfinder/internal/bootstrap"
	"github.com/woningfinder/woningfinder/internal/corporation"
	"github.com/woningfinder/woningfinder/pkg/config"
	"github.com/woningfinder/woningfinder/pkg/logging"
	"go.uber.org/zap"

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

// workerMode defines if housing-finder should be run with internal cron (always running and cron managed by go)
// or via command line using `housing-finder standalone`
var workerMode = true

func main() {
	logger := logging.NewZapLoggerWithSentry(config.MustGetString("SENTRY_DSN"))

	// check if should be run is workerMode
	workerMode = !(len(os.Args) > 1 && os.Args[1] == "standalone")

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
	if workerMode {
		// get time location
		nld, err := time.LoadLocation("Europe/Amsterdam")
		if err != nil {
			logger.Sugar().Fatal(err)
		}

		// instantiate cron
		c := cron.New(cron.WithLocation(nld), cron.WithSeconds(), cron.WithLogger(cron.VerbosePrintfLogger(log.New(os.Stdout, "cron: ", log.LstdFlags))))
		parser := cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)

		// populate crons
		for _, spec := range config.MustGetStringList("HOUSING_FINDER_SCHEDULE") {
			_, err := parser.Parse(spec)
			if err != nil {
				logger.Sugar().Errorf("error when parsing cron spec %s: %w", spec, err)
				continue
			}
			c.AddFunc(spec, func() {
				findHousing(logger, clientProvider, corporationService)
			})
		}

		// start cron scheduler
		c.Run()
	} else {
		findHousing(logger, clientProvider, corporationService)
	}
}

func findHousing(logger *zap.Logger, clientProvider corporation.ClientProvider, corporationService corporation.Service) {
	wg := sync.WaitGroup{}
	for _, c := range *clientProvider.List() {
		client, err := clientProvider.Get(c)
		if err != nil {
			logger.Sugar().Warn(err)
			continue
		}
		wg.Add(1)
		go func(corporation corporation.Corporation, client corporation.Client) {
			defer wg.Done()

			if err := corporationService.PublishOffers(client, corporation); err != nil {
				logger.Sugar().Error(err)
			}
		}(c, client)
	}
	wg.Wait()
}
