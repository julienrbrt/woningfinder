package corporation

import (
	"github.com/woningfinder/woningfinder/internal/database"
	"github.com/woningfinder/woningfinder/internal/domain/entity"
	"github.com/woningfinder/woningfinder/pkg/logging"
)

// Service permits to handle the persistence of a corporation
type Service interface {
	// Corporation
	CreateOrUpdateCorporation(corporation entity.Corporation) error
	GetCorporation(name string) (*entity.Corporation, error)
	DeleteCorporation(corp entity.Corporation) error

	// Cities
	AddCities(cities []entity.City, corporation entity.Corporation) error
	GetCities() (*[]entity.City, error)
	GetCity(name string) (*entity.City, error)
	DeleteCity(city entity.City) error
}

type service struct {
	logger   *logging.Logger
	dbClient database.DBClient
}

// NewService instantiate the corporation service
func NewService(logger *logging.Logger, dbClient database.DBClient) Service {
	return &service{
		logger:   logger,
		dbClient: dbClient,
	}
}
