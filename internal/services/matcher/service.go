package matcher

import (
	"context"

	"github.com/julienrbrt/woningfinder/internal/corporation"
	"github.com/julienrbrt/woningfinder/internal/corporation/connector"
	"github.com/julienrbrt/woningfinder/internal/customer/matcher"
	"github.com/julienrbrt/woningfinder/internal/database"
	corporationService "github.com/julienrbrt/woningfinder/internal/services/corporation"
	emailService "github.com/julienrbrt/woningfinder/internal/services/email"
	userService "github.com/julienrbrt/woningfinder/internal/services/user"
	"github.com/julienrbrt/woningfinder/pkg/digitalocean/spaces"
	"github.com/julienrbrt/woningfinder/pkg/logging"
)

// offersQueue defines on which queue the corporation offers are sent via redis
const offersQueue = "queue:offers"

// Service permits to handle the persistence of matcher
type Service interface {
	SendOffers(offers corporation.Offers) error
	RetrieveOffers(ch chan<- corporation.Offers) error

	MatchOffer(ctx context.Context, offers corporation.Offers) error
}

type service struct {
	logger             *logging.Logger
	redisClient        database.RedisClient
	userService        userService.Service
	emailService       emailService.Service
	corporationService corporationService.Service
	spacesClient       spaces.Client
	matcher            matcher.Matcher
	connectorProvider  connector.ConnectorProvider
}

// NewService instantiate the matcher service
func NewService(logger *logging.Logger, redisClient database.RedisClient, userService userService.Service, emailService emailService.Service, corporationService corporationService.Service, spacesClient spaces.Client, matcher matcher.Matcher, connectorProvider connector.ConnectorProvider) Service {
	return &service{
		logger:             logger,
		redisClient:        redisClient,
		userService:        userService,
		emailService:       emailService,
		corporationService: corporationService,
		spacesClient:       spacesClient,
		matcher:            matcher,
		connectorProvider:  connectorProvider,
	}
}
