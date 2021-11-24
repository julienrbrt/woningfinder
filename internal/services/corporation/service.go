package corporation

import (
	"github.com/julienrbrt/woningfinder/internal/corporation"
	"github.com/julienrbrt/woningfinder/internal/corporation/city"
	"github.com/julienrbrt/woningfinder/internal/database"
	"github.com/julienrbrt/woningfinder/pkg/logging"
)

// Service permits to handle the persistence of a corporation
type Service interface {
	// Corporation
	CreateOrUpdateCorporation(corporation corporation.Corporation) error
	GetCorporation(name string) (*corporation.Corporation, error)

	// Cities
	LinkCities(cities []city.City, hasLocation bool, corp ...corporation.Corporation) error
	GetCities() ([]*city.City, error)
	GetCity(name string) (*city.City, error)
}

type service struct {
	logger    *logging.Logger
	dbClient  database.DBClient
	suggester city.Suggester
}

// NewService instantiate the corporation service
func NewService(logger *logging.Logger, dbClient database.DBClient, suggester city.Suggester) Service {
	return &service{
		logger:    logger,
		dbClient:  dbClient,
		suggester: suggester,
	}
}
