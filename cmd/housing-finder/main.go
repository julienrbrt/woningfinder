package main

import (
	"os"
	"time"

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
	logger := logging.NewZapLoggerWithoutSentry()
	jwtAuth := auth.CreateJWTAuthenticationToken(config.MustGetString("JWT_SECRET"))

	dbClient := bootstrap.CreateDBClient(logger)
	redisClient := bootstrap.CreateRedisClient(logger)
	mapboxClient := bootstrap.CreateMapboxClient(logger, redisClient)
	spacesClient := bootstrap.CreateDOSpacesClient(logger)
	emailClient := bootstrap.CreateEmailClient()

	clientProvider := bootstrapCorporation.CreateClientProvider(logger, mapboxClient)
	corporationService := corporationService.NewService(logger, dbClient)
	userService := userService.NewService(logger, dbClient, config.MustGetString("AES_SECRET"), clientProvider, corporationService)
	emailService := emailService.NewService(logger, emailClient, jwtAuth)
	matcherService := matcherService.NewService(logger, redisClient, userService, emailService, corporationService, spacesClient, matcher.NewMatcher(), clientProvider)

	if len(os.Args) != 2 {
		logger.Fatal("usage: housing-finder <corporation-name>")
	}

	// get corporation client
	corp, err := clientProvider.GetCorporation(os.Args[1])
	if err != nil {
		logger.Fatal("error while getting corporation", zap.String("got", os.Args[1]), zap.Error(err))
	}

	client, err := clientProvider.Get(os.Args[1])
	if err != nil {
		logger.Fatal("error while getting corporation client", zap.String("got", os.Args[1]), zap.Error(err))
	}

	logger.Info("housing-finder started", zap.String("corporation", corp.Name))

	ch := make(chan corporation.Offer)
	go func(ch chan corporation.Offer) {
		defer close(ch)

		if err := client.FetchOffers(ch); err != nil {
			logger.Error("error while fetching offers", zap.String("corporation", corp.Name), zap.Error(err))
		}
		defer close(ch)
	}(ch)

	offers := corporation.Offers{
		CorporationName: corp.Name,
		Offer:           []corporation.Offer{},
	}

	// batch send offers every 5 seconds
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	counter := 0
	for {
		select {
		case <-ticker.C:
			if len(offers.Offer) == 0 {
				continue
			}

			logger.Info("housing-finder job sending offers", zap.String("corporation", corp.Name), zap.Int("offers", len(offers.Offer)))

			if err := matcherService.SendOffers(offers); err != nil {
				logger.Error("error while sending offer to redis queue", zap.String("corporation", offers.CorporationName), zap.Error(err))
			}

			counter += len(offers.Offer)
			offers.Offer = []corporation.Offer{}
		case offer, ok := <-ch:
			if ok {
				logger.Info("new offer parsed...", zap.Any("offer", offer))
				offers.Offer = append(offers.Offer, offer)
			}

			if !ok && len(offers.Offer) == 0 {
				logger.Info("housing-finder finished", zap.Int("offers sent", counter), zap.String("corporation", corp.Name))
				return
			}
		}
	}
}
