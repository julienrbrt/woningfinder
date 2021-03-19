package matcher

import (
	"context"

	"github.com/woningfinder/woningfinder/internal/corporation"
	"github.com/woningfinder/woningfinder/internal/database"
	"github.com/woningfinder/woningfinder/internal/domain/entity"
	corporationService "github.com/woningfinder/woningfinder/internal/services/corporation"
	userService "github.com/woningfinder/woningfinder/internal/services/user"
	"github.com/woningfinder/woningfinder/pkg/logging"
)

// offersQueue defines on which queue the corporation offers are sent via redis
const offersQueue = "queue:offers"

// Service permits to handle the persistence of matcher
type Service interface {
	PublishOffers(client corporation.Client, corporation entity.Corporation) error
	SubscribeOffers(ch chan<- entity.OfferList) error

	MatchOffer(ctx context.Context, offers entity.OfferList) error
}

type service struct {
	logger             *logging.Logger
	redisClient        database.RedisClient
	userService        userService.Service
	corporationService corporationService.Service
	clientProvider     corporation.ClientProvider
}

// NewService instantiate the matcher service
func NewService(logger *logging.Logger, redisClient database.RedisClient, userService userService.Service, corporationService corporationService.Service, clientProvider corporation.ClientProvider) Service {
	return &service{
		logger:             logger,
		redisClient:        redisClient,
		userService:        userService,
		corporationService: corporationService,
		clientProvider:     clientProvider,
	}
}
