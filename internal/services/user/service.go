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
	CreateUser(user *customer.User) error
	GetUser(email string) (*customer.User, error)
	DeleteUser(email string) error
	ConfirmUser(email string) error
	SetStripeCustomerID(user *customer.User, stripeID string) error
	ConfirmSubscription(stripeID string) error
	UpdateSubscriptionStatus(stripeID string, status bool) error
	UpdateUser(user *customer.User) error

	GetUsersWithGivenCorporationCredentials(corporationName string) ([]*customer.User, error)
	GetWeeklyUpdateUsers() ([]*customer.User, error)

	// Housing Preferences
	CreateHousingPreferences(userID uint, preferences *customer.HousingPreferences) error
	GetHousingPreferences(userID uint) (customer.HousingPreferences, error)
	UpdateHousingPreferences(userID uint, preferences *customer.HousingPreferences) error
	DeleteHousingPreferences(userID uint) error
	GetHousingPreferencesMatchingCorporation(userID uint) ([]*corporation.Corporation, error)
	CreateHousingPreferencesMatch(userID uint, offer corporation.Offer, corporationName, pictureURL string) error

	// Corporation Credentials
	HasCorporationCredentials(userID uint) (bool, error)
	CreateCorporationCredentials(userID uint, credentials *customer.CorporationCredentials) error
	GetCorporationCredentials(userID uint, corporationName string) (*customer.CorporationCredentials, error)
	DeleteCorporationCredentials(userID uint, corporationName string) error
	UpdateCorporationCredentialsFailureCount(userID uint, corporationName string, failureCount int) error
	DecryptCredentials(credentials *customer.CorporationCredentials) (*customer.CorporationCredentials, error)

	// Waiting list
	CreateWaitingList(w *customer.WaitingList) error
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
