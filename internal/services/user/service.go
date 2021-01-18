package user

import (
	"github.com/woningfinder/woningfinder/internal/corporation"
	"github.com/woningfinder/woningfinder/internal/database"
	"github.com/woningfinder/woningfinder/internal/domain/entity"
	corporationService "github.com/woningfinder/woningfinder/internal/services/corporation"
	"github.com/woningfinder/woningfinder/pkg/logging"
)

// UserService permits to handle the persistence of an user
type UserService interface {
	CreateUser(u entity.User) error
	GetUser(email string) (*entity.User, error)
	DeleteUser(u *entity.User) error

	CreateHousingPreferences(u *entity.User, preferences []entity.HousingPreferences) error
	GetHousingPreferences(u *entity.User) ([]entity.HousingPreferences, error)
	DeleteHousingPreferences(u *entity.User) error
	CreateHousingPreferencesMatch(u *entity.User, offer entity.Offer, corporationName string) error

	CreateCorporationCredentials(u *entity.User, credentials entity.CorporationCredentials) error
	GetCorporationCredentials(u *entity.User, corporation entity.Corporation) (*entity.CorporationCredentials, error)
	GetAllCorporationCredentials(corporation entity.Corporation) ([]entity.CorporationCredentials, error)
	DeleteCorporationCredentials(u *entity.User, corporation entity.Corporation) error

	MatchOffer(offers entity.OfferList) error
}

type service struct {
	logger             *logging.Logger
	dbClient           database.DBClient
	redisClient        database.RedisClient
	aesSecret          string
	clientProvider     corporation.ClientProvider
	corporationService corporationService.CorporationService
}

func NewService(logger *logging.Logger, dbClient database.DBClient, redisClient database.RedisClient, aesSecret string, clientProvider corporation.ClientProvider, corporationService corporationService.CorporationService) UserService {
	return &service{
		logger:             logger,
		dbClient:           dbClient,
		redisClient:        redisClient,
		aesSecret:          aesSecret,
		clientProvider:     clientProvider,
		corporationService: corporationService,
	}
}
