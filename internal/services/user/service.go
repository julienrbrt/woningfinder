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
	ConfirmUser(email string) (*customer.User, error)
	ConfirmPayment(email string) (*customer.User, error)
	GetWeeklyUpdateUsers() ([]*customer.User, error)

	// Waiting List
	CreateWaitingList(w *customer.WaitingList) error

	// Housing Preferences
	CreateHousingPreferences(user *customer.User, preferences customer.HousingPreferences) error
	GetHousingPreferences(user *customer.User) (customer.HousingPreferences, error)
	UpdateHousingPreferences(user *customer.User, housingPreferences customer.HousingPreferences) error
	DeleteHousingPreferences(user *customer.User) error
	GetHousingPreferencesMatchingCorporation(user *customer.User) ([]corporation.Corporation, error)
	CreateHousingPreferencesMatch(user *customer.User, offer corporation.Offer, corporationName string) error

	// Corporation Credentials
	CreateCorporationCredentials(userID uint, credentials customer.CorporationCredentials) error
	GetCorporationCredentials(userID uint, corporationName string) (*customer.CorporationCredentials, error)
	GetAllCorporationCredentials(corporationName string) ([]customer.CorporationCredentials, error)
	DeleteCorporationCredentials(userID uint, corporationName string) error
	DecryptCredentials(credentials *customer.CorporationCredentials) (*customer.CorporationCredentials, error)
	UpdateCorporationCredentialsFailureCount(userID uint, corporationName string, failureCount int) error
}

type service struct {
	logger             *logging.Logger
	dbClient           database.DBClient
	aesSecret          string
	clientProvider     connector.ClientProvider
	corporationService corporationService.Service
}

// NewService instantiate the user service
func NewService(logger *logging.Logger, dbClient database.DBClient, aesSecret string, clientProvider connector.ClientProvider, corporationService corporationService.Service) Service {
	return &service{
		logger:             logger,
		dbClient:           dbClient,
		aesSecret:          aesSecret,
		clientProvider:     clientProvider,
		corporationService: corporationService,
	}
}
