package corporation

import (
	"github.com/woningfinder/woningfinder/internal/corporation"
	"github.com/woningfinder/woningfinder/internal/database"
	"github.com/woningfinder/woningfinder/pkg/logging"
)

// Service permits to handle the persistence of a corporation
type Service interface {
	// Corporation
	CreateOrUpdateCorporation(corporation corporation.Corporation) error
	GetCorporation(name string) (*corporation.Corporation, error)

	// Cities
	LinkCities(cities []corporation.City, corp ...corporation.Corporation) error
	GetCities() ([]corporation.City, error)
	GetCity(name string) (corporation.City, error)
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
