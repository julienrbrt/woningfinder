package user

import (
	"github.com/woningfinder/woningfinder/internal/corporation"
	"github.com/woningfinder/woningfinder/internal/corporation/connector"
	"github.com/woningfinder/woningfinder/internal/customer"
	"github.com/woningfinder/woningfinder/internal/database"
	corporationService "github.com/woningfinder/woningfinder/internal/services/corporation"
	"github.com/woningfinder/woningfinder/pkg/logging"
)

// Service permits to handle the persistence of an user
type Service interface {
	// Users
	CreateUser(u *customer.User) error
	GetUser(search *customer.User) (*customer.User, error)
	DeleteUser(u *customer.User) error
	GetWeeklyUpdateUsers() ([]*customer.User, error)

	// Payment
	SetPaid(userID uint, plan customer.Plan) error

	// Housing Preferences
	CreateHousingPreferences(u *customer.User, preferences customer.HousingPreferences) error
	GetHousingPreferences(u *customer.User) (customer.HousingPreferences, error)
	CreateHousingPreferencesMatch(u *customer.User, offer corporation.Offer, corporationName string) error
	GetHousingPreferencesMatchingCorporation(u *customer.User) ([]corporation.Corporation, error)

	// Corporation Credentials
	CreateCorporationCredentials(userID uint, credentials customer.CorporationCredentials) error
	GetCorporationCredentials(userID uint, corporationName string) (*customer.CorporationCredentials, error)
	GetAllCorporationCredentials(corporationName string) ([]customer.CorporationCredentials, error)
	DeleteCorporationCredentials(userID uint, corporationName string) error
	DecryptCredentials(credentials *customer.CorporationCredentials) (*customer.CorporationCredentials, error)
}

type service struct {
	logger             *logging.Logger
	dbClient           database.DBClient
	redisClient        database.RedisClient
	aesSecret          string
	clientProvider     connector.ClientProvider
	corporationService corporationService.Service
}

// NewService instantiate the user service
func NewService(logger *logging.Logger, dbClient database.DBClient, redisClient database.RedisClient, aesSecret string, clientProvider connector.ClientProvider, corporationService corporationService.Service) Service {
	return &service{
		logger:             logger,
		dbClient:           dbClient,
		redisClient:        redisClient,
		aesSecret:          aesSecret,
		clientProvider:     clientProvider,
		corporationService: corporationService,
	}
}
