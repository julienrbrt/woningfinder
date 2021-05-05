package main

import (
	"context"

	"github.com/joho/godotenv"
	"github.com/woningfinder/woningfinder/internal/auth"
	"github.com/woningfinder/woningfinder/internal/bootstrap"
	bootstrapCorporation "github.com/woningfinder/woningfinder/internal/bootstrap/corporation"
	"github.com/woningfinder/woningfinder/internal/corporation"
	"github.com/woningfinder/woningfinder/internal/customer/matcher"
	corporationService "github.com/woningfinder/woningfinder/internal/services/corporation"
	matcherService "github.com/woningfinder/woningfinder/internal/services/matcher"
	notificationService "github.com/woningfinder/woningfinder/internal/services/notification"
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
}

func main() {
	logger := logging.NewZapLogger(config.GetBoolOrDefault("APP_DEBUG", false), config.MustGetString("SENTRY_DSN"))
	jwtAuth := auth.CreateJWTAuthenticationToken(config.MustGetString("JWT_SECRET"))

	dbClient := bootstrap.CreateDBClient(logger)
	redisClient := bootstrap.CreateRedisClient(logger)
	emailClient := bootstrap.CreateEmailClient()

	clientProvider := bootstrapCorporation.CreateClientProvider(logger, nil) // mapboxClient not required in the matcher
	corporationService := corporationService.NewService(logger, dbClient)
	userService := userService.NewService(logger, dbClient, redisClient, config.MustGetString("AES_SECRET"), clientProvider, corporationService)
	notificationService := notificationService.NewService(logger, emailClient, jwtAuth)
	matcherService := matcherService.NewService(logger, redisClient, userService, notificationService, corporationService, matcher.NewMatcher(), clientProvider)

	// subscribe to offers queue inside a new go routine
	ch := make(chan corporation.Offers)
	go func(ch chan corporation.Offers) {
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
