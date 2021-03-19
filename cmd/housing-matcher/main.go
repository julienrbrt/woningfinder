package main

import (
	"context"

	"github.com/joho/godotenv"
	"github.com/woningfinder/woningfinder/internal/bootstrap"
	"github.com/woningfinder/woningfinder/internal/domain/entity"
	corporationService "github.com/woningfinder/woningfinder/internal/services/corporation"
	matcherService "github.com/woningfinder/woningfinder/internal/services/matcher"
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

	dbClient := bootstrap.CreateDBClient(logger)
	redisClient := bootstrap.CreateRedisClient(logger)
	clientProvider := bootstrap.CreateClientProvider(logger, nil) // mapboxClient not required in the matcher
	corporationService := corporationService.NewService(logger, dbClient)
	userService := userService.NewService(logger, dbClient, redisClient, config.MustGetString("AES_SECRET"), clientProvider, corporationService)
	matcherService := matcherService.NewService(logger, redisClient, userService, corporationService, clientProvider)

	// subscribe to offers queue inside a new go routine
	ch := make(chan entity.OfferList)
	go func(ch chan entity.OfferList) {
		if err := matcherService.SubscribeOffers(ch); err != nil {
			logger.Sugar().Fatal(err)
		}
	}(ch)

	// match offer
	for offers := range ch {
		if err := matcherService.MatchOffer(context.Background(), offers); err != nil {
			logger.Sugar().Errorf("error while maching offers for corporation %s: %w", offers.Corporation.Name, err)
		}
	}
}
