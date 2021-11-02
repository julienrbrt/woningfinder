package main

import (
	"github.com/joho/godotenv"
	"github.com/woningfinder/woningfinder/internal/bootstrap"
	"github.com/woningfinder/woningfinder/internal/corporation/city"
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
	dbClient := bootstrap.CreateDBClient(logger)

	var cities []city.City
	if err := dbClient.Conn().Model(&cities).Where("latitude IS NULL OR longitude IS NULL").Select(); err != nil {
		logger.Fatal("error while getting cities without location", zap.Error(err))
	}

	for _, city := range cities {
		logger.Info("updating city", zap.String("city", city.Name))
		if _, err := dbClient.Conn().Model(&city).OnConflict("(name) DO UPDATE").Insert(); err != nil {
			logger.Error("failed updating city", zap.String("city", city.Name), zap.Error(err))
		}
	}
}
