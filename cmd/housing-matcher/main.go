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
	emailService "github.com/woningfinder/woningfinder/internal/services/email"
	matcherService "github.com/woningfinder/woningfinder/internal/services/matcher"
	userService "github.com/woningfinder/woningfinder/internal/services/user"
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
	jwtAuth := auth.CreateJWTAuthenticationToken(config.MustGetString("JWT_SECRET"))

	dbClient := bootstrap.CreateDBClient(logger)
	redisClient := bootstrap.CreateRedisClient(logger)
	emailClient := bootstrap.CreateEmailClient()
	spacesClient := bootstrap.CreateDOSpacesClient(logger)

	clientProvider := bootstrapCorporation.CreateClientProvider(logger, nil) // mapboxClient not required in the matcher
	corporationService := corporationService.NewService(logger, dbClient)
	userService := userService.NewService(logger, dbClient, config.MustGetString("AES_SECRET"), clientProvider, corporationService)
	emailService := emailService.NewService(logger, emailClient, jwtAuth)
	matcherService := matcherService.NewService(logger, redisClient, userService, emailService, corporationService, spacesClient, matcher.NewMatcher(), clientProvider)

	// subscribe to offers queue inside a new go routine
	ch := make(chan corporation.Offers)
	go func(ch chan corporation.Offers) {
		if err := matcherService.RetrieveOffers(ch); err != nil {
			logger.Fatal("failed subscribing to offers", zap.Error(err))
		}
	}(ch)

	// match offer
	for offers := range ch {
		logger.Info("received offers", zap.String("corporation", offers.CorporationName), zap.Int("count", len(offers.Offer)))

		if err := matcherService.MatchOffer(context.Background(), offers); err != nil {
			logger.Error("error while maching offers", zap.String("corporation", offers.CorporationName), zap.Error(err))
		}
	}
}
