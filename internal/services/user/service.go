package user

import (
	"go.uber.org/zap"

	"github.com/woningfinder/woningfinder/internal/corporation"
	"github.com/woningfinder/woningfinder/internal/database"
	"github.com/woningfinder/woningfinder/internal/domain/entity"
	corporationService "github.com/woningfinder/woningfinder/internal/services/corporation"
)

// Service permits to handle the persistence of an user
type Service interface {
	CreateUser(u *entity.User) (*entity.User, error)
	GetUser(email string) (*entity.User, error)
	DeleteUser(u *entity.User) error

	CreateHousingPreferences(u *entity.User, pref []entity.HousingPreferences) error
	GetHousingPreferences(u *entity.User) (*[]entity.HousingPreferences, error)
	DeleteHousingPreferences(u *entity.User) error

	CreateCorporationCredentials(u *entity.User, credentials entity.CorporationCredentials) error
	GetCorporationCredentials(u *entity.User, corporation entity.Corporation) (*entity.CorporationCredentials, error)
	DeleteCorporationCredentials(u *entity.User, corporation entity.Corporation) error

	MatchOffer(offers entity.OfferList) error
}

type service struct {
	logger             *zap.Logger
	dbClient           database.DBClient
	redisClient        database.RedisClient
	aesSecret          string
	clientProvider     corporation.ClientProvider
	corporationService corporationService.Service
}

func NewService(logger *zap.Logger, dbClient database.DBClient, redisClient database.RedisClient, aesSecret string, clientProvider corporation.ClientProvider, corporationService corporationService.Service) Service {
	return &service{
		logger:             logger,
		dbClient:           dbClient,
		redisClient:        redisClient,
		aesSecret:          aesSecret,
		clientProvider:     clientProvider,
		corporationService: corporationService,
	}
}
