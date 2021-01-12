package corporation

import (
	"github.com/woningfinder/woningfinder/internal/corporation"
	"github.com/woningfinder/woningfinder/internal/database"
	"github.com/woningfinder/woningfinder/internal/domain/entity"
	"go.uber.org/zap"
)

// Service permits to handle the persistence of a corporation
type Service interface {
	CreateOrUpdateCorporation(corporation *[]entity.Corporation) (*[]entity.Corporation, error)
	GetCity(name string) (*entity.City, error)

	PublishOffers(client corporation.Client, corporation entity.Corporation) error
	SubscribeOffers(offerCh chan<- entity.OfferList)
}

type service struct {
	logger      *zap.Logger
	dbClient    database.DBClient
	redisClient database.RedisClient
}

func NewService(logger *zap.Logger, dbClient database.DBClient, redisClient database.RedisClient) Service {
	return &service{
		logger:      logger,
		dbClient:    dbClient,
		redisClient: redisClient,
	}
}
