package main

import (
	"github.com/joho/godotenv"
	"github.com/woningfinder/woningfinder/internal/bootstrap"
	"github.com/woningfinder/woningfinder/internal/customer/matcher"
	"github.com/woningfinder/woningfinder/internal/services/corporation"
	matcherService "github.com/woningfinder/woningfinder/internal/services/matcher"
	"github.com/woningfinder/woningfinder/internal/services/user"
	"github.com/woningfinder/woningfinder/pkg/config"
	"github.com/woningfinder/woningfinder/pkg/logging"
)

// init is invoked before main()
func init() {
	// loads values from .env into the system
	// fallback to system env if unexisting
	// if not defined on system, panics
	if err := godotenv.Load("../../../.env"); err != nil {
		_ = config.MustGetString("APP_NAME")
	}
}

func main() {
	logger := logging.NewZapLoggerWithoutSentry()

	dbClient := bootstrap.CreateDBClient(logger)
	redisClient := bootstrap.CreateRedisClient(logger)
	mapboxClient := bootstrap.CreateMapboxClient()

	clientProvider := bootstrap.CreateClientProvider(logger, mapboxClient)
	corporationService := corporation.NewService(logger, dbClient)
	userService := user.NewService(logger, dbClient, redisClient, config.MustGetString("AES_SECRET"), clientProvider, corporationService)
	matcherService := matcherService.NewService(logger, redisClient, userService, nil, corporationService, matcher.NewMatcher(), clientProvider)

	// populate crons
	for _, corp := range clientProvider.List() {
		corp := corp // https://github.com/golang/go/wiki/CommonMistakes#using-reference-to-loop-iterator-variable

		// get corporation client
		client, err := clientProvider.Get(corp.Name)
		if err != nil {
			logger.Sugar().Error(err)
			continue
		}

		if err := matcherService.PublishOffers(client, corp); err != nil {
			logger.Sugar().Error(err)
		}
	}
}
