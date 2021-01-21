package user

import (
	"github.com/woningfinder/woningfinder/internal/corporation"
	"github.com/woningfinder/woningfinder/internal/database"
	"github.com/woningfinder/woningfinder/internal/domain/entity"
	corporationService "github.com/woningfinder/woningfinder/internal/services/corporation"
	"github.com/woningfinder/woningfinder/pkg/logging"
)

// Service permits to handle the persistence of an user
type Service interface {
	// Users
	CreateUser(u *entity.User) error
	GetUser(search *entity.User) (*entity.User, error)
	DeleteUser(u *entity.User) error

	// Payment
	SetPaid(u *entity.User, plan entity.Plan) error
	SetExpired(u *entity.User) error

	// Housing Preferences
	CreateHousingPreferences(u *entity.User, preferences []entity.HousingPreferences) error
	GetHousingPreferences(u *entity.User) ([]entity.HousingPreferences, error)
	DeleteHousingPreferences(u *entity.User) error
	GetHousingPreferencesMatchingCorporation(u *entity.User) ([]entity.Corporation, error)

	CreateHousingPreferencesMatch(u *entity.User, offer entity.Offer, corporationName string) error

	// Corporation Credentials
	CreateCorporationCredentials(u *entity.User, credentials entity.CorporationCredentials) error
	GetCorporationCredentials(u *entity.User, corporation entity.Corporation) (*entity.CorporationCredentials, error)
	GetAllCorporationCredentials(corporation entity.Corporation) ([]entity.CorporationCredentials, error)
	DeleteCorporationCredentials(u *entity.User, corporation entity.Corporation) error
	ValidateCredentials(credentials entity.CorporationCredentials) error
	DecryptCredentials(credentials *entity.CorporationCredentials) (*entity.CorporationCredentials, error)
}

type service struct {
	logger             *logging.Logger
	dbClient           database.DBClient
	redisClient        database.RedisClient
	aesSecret          string
	clientProvider     corporation.ClientProvider
	corporationService corporationService.Service
}

// NewService instantiate the user service
func NewService(logger *logging.Logger, dbClient database.DBClient, redisClient database.RedisClient, aesSecret string, clientProvider corporation.ClientProvider, corporationService corporationService.Service) Service {
	return &service{
		logger:             logger,
		dbClient:           dbClient,
		redisClient:        redisClient,
		aesSecret:          aesSecret,
		clientProvider:     clientProvider,
		corporationService: corporationService,
	}
}
