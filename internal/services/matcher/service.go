package matcher

import (
	"context"

	"github.com/woningfinder/woningfinder/internal/corporation"
	"github.com/woningfinder/woningfinder/internal/corporation/connector"
	"github.com/woningfinder/woningfinder/internal/customer/matcher"
	"github.com/woningfinder/woningfinder/internal/database"
	corporationService "github.com/woningfinder/woningfinder/internal/services/corporation"
	emailService "github.com/woningfinder/woningfinder/internal/services/email"
	userService "github.com/woningfinder/woningfinder/internal/services/user"
	"github.com/woningfinder/woningfinder/pkg/logging"
)

// offersQueue defines on which queue the corporation offers are sent via redis
const offersQueue = "queue:offers"

// Service permits to handle the persistence of matcher
type Service interface {
	PushOffers(client connector.Client, corporation corporation.Corporation) error
	SubscribeOffers(ch chan<- corporation.Offers) error

	MatchOffer(ctx context.Context, offers corporation.Offers) error
}

type service struct {
	logger             *logging.Logger
	redisClient        database.RedisClient
	userService        userService.Service
	emailService       emailService.Service
	corporationService corporationService.Service
	matcher            matcher.Matcher
	clientProvider     connector.ClientProvider
}

// NewService instantiate the matcher service
func NewService(logger *logging.Logger, redisClient database.RedisClient, userService userService.Service, emailService emailService.Service, corporationService corporationService.Service, matcher matcher.Matcher, clientProvider connector.ClientProvider) Service {
	return &service{
		logger:             logger,
		redisClient:        redisClient,
		userService:        userService,
		emailService:       emailService,
		corporationService: corporationService,
		matcher:            matcher,
		clientProvider:     clientProvider,
	}
}
