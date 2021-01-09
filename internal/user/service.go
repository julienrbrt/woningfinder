package user

import (
	"go.uber.org/zap"

	"github.com/woningfinder/woningfinder/internal/corporation"
	"github.com/woningfinder/woningfinder/internal/database"
)

// Service permits to handle the persistence of an user
type Service interface {
	CreateUser(u *User) (*User, error)
	GetUser(email string) (*User, error)
	DeleteUser(u *User) error

	CreateHousingPreferences(u *User, pref []HousingPreferences) error
	GetHousingPreferences(u *User) (*[]HousingPreferences, error)
	DeleteHousingPreferences(u *User) error

	CreateCorporationCredentials(u *User, credentials CorporationCredentials) error
	GetCorporationCredentials(u *User, corporation corporation.Corporation) (*CorporationCredentials, error)
	DeleteCorporationCredentials(u *User, corporation corporation.Corporation) error

	MatchOffer(offers corporation.OfferList) error
}

type service struct {
	logger             *zap.Logger
	dbClient           database.DBClient
	redisClient        database.RedisClient
	aesSecret          string
	clientProvider     corporation.ClientProvider
	corporationService corporation.Service
}

func NewService(logger *zap.Logger, dbClient database.DBClient, redisClient database.RedisClient, aesSecret string, clientProvider corporation.ClientProvider, corporationService corporation.Service) Service {
	return &service{
		logger:             logger,
		dbClient:           dbClient,
		redisClient:        redisClient,
		aesSecret:          aesSecret,
		clientProvider:     clientProvider,
		corporationService: corporationService,
	}
}
