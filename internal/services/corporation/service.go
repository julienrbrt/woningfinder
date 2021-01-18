package corporation

import (
	"github.com/woningfinder/woningfinder/internal/corporation"
	"github.com/woningfinder/woningfinder/internal/database"
	"github.com/woningfinder/woningfinder/internal/domain/entity"
	"github.com/woningfinder/woningfinder/pkg/logging"
)

// CorporationService permits to handle the persistence of a corporation
type CorporationService interface {
	CreateOrUpdateCorporation(corporation entity.Corporation) error
	GetCorporation(name string) (*entity.Corporation, error)
	DeleteCorporation(corp entity.Corporation) error
	// GetCorporationByCity(city string) ([]entity.Corporation, error) // TODO for login page

	AddCities(city []entity.City) ([]entity.City, error)
	GetCities() (*[]entity.City, error)
	GetCity(name string) (*entity.City, error)
	DeleteCity(city entity.City) error

	PublishOffers(client corporation.Client, corporation entity.Corporation) error
	SubscribeOffers(offerCh chan<- entity.OfferList) error
}

type service struct {
	logger      *logging.Logger
	dbClient    database.DBClient
	redisClient database.RedisClient
}

func NewService(logger *logging.Logger, dbClient database.DBClient, redisClient database.RedisClient) CorporationService {
	return &service{
		logger:      logger,
		dbClient:    dbClient,
		redisClient: redisClient,
	}
}
