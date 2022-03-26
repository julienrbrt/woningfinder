package main

import (
	"github.com/joho/godotenv"
	"github.com/julienrbrt/woningfinder/internal/bootstrap"
	bootstrapCorporation "github.com/julienrbrt/woningfinder/internal/bootstrap/corporation"
	"github.com/julienrbrt/woningfinder/internal/corporation/city"
	"github.com/julienrbrt/woningfinder/internal/services/corporation"
	"github.com/julienrbrt/woningfinder/pkg/config"
	"github.com/julienrbrt/woningfinder/pkg/logging"
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
	dbClient := bootstrap.CreateDBClient(logger)
	connectorProvider := bootstrapCorporation.CreateConnectorProvider(logger, nil)
	cityTable := connectorProvider.GetCities()
	corporationService := corporation.NewService(logger, dbClient, city.NewSuggester(cityTable))

	var cities []city.City
	if err := dbClient.Conn().Model(&cities).Where("latitude IS NULL OR longitude IS NULL").Select(); err != nil {
		logger.Fatal("error while getting cities without location", zap.Error(err))
	}

	for _, c := range cities {
		city, ok := cityTable[c.Name]
		if !ok {
			logger.Warn("failed finding city", zap.String("city", c.Name))
			continue
		}

		logger.Info("updating city", zap.String("city", city.Name))
		if _, err := dbClient.Conn().Model(&city).OnConflict("(name) DO UPDATE").Insert(); err != nil {
			logger.Error("failed updating city", zap.String("city", city.Name), zap.Error(err))
		}
	}

	// update corporations cities
	for _, corp := range connectorProvider.GetCorporations() {
		if err := corporationService.LinkCities(corp.Cities, true, corp); err != nil {
			logger.Error("failed updating corporation", zap.String("corporation", corp.Name), zap.Error(err))
		}
	}
}
