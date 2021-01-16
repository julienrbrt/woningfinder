package main

import (
	"github.com/woningfinder/woningfinder/internal/domain/entity"
	"github.com/woningfinder/woningfinder/pkg/logging"

	"github.com/joho/godotenv"
	"github.com/woningfinder/woningfinder/internal/bootstrap"
	"github.com/woningfinder/woningfinder/internal/services/corporation"
	"github.com/woningfinder/woningfinder/internal/services/user"
	"github.com/woningfinder/woningfinder/pkg/config"
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

	dbClient := bootstrap.CreateDBClient(logger)
	redisClient := bootstrap.CreateRedisClient(logger)
	clientProvider := bootstrap.CreateClientProvider(logger, nil)
	corporationService := corporation.NewService(logger, dbClient, redisClient)
	userService := user.NewService(logger, dbClient, redisClient, config.MustGetString("AES_SECRET"), clientProvider, corporationService)

	offerList := make(chan entity.OfferList)
	// subscribe to pub/sub messages inside a new goroutine
	go func(offerList chan entity.OfferList) {
		if err := corporationService.SubscribeOffers(offerList); err != nil {
			logger.Sugar().Fatal(err)
		}
	}(offerList)

	for o := range offerList {
		if err := userService.MatchOffer(o); err != nil {
			logger.Sugar().Errorf("error while maching offers for corporation %s: %w", o.Corporation.Name, err)
		}
	}

}
