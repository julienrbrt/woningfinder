package corporation_test

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/go-redis/redis"

	"github.com/stretchr/testify/assert"

	"github.com/woningfinder/woningfinder/internal/domain/entity"

	"github.com/woningfinder/woningfinder/internal/corporation"
	"github.com/woningfinder/woningfinder/internal/corporation/dewoonplaats"

	"github.com/woningfinder/woningfinder/internal/database"
	corpServ "github.com/woningfinder/woningfinder/internal/services/corporation"

	"github.com/woningfinder/woningfinder/pkg/logging"
)

func Test_PublishOffers_CorporationClientError(t *testing.T) {
	a := assert.New(t)

	err := errors.New("foo")
	logger := logging.NewZapLoggerWithoutSentry()
	redisMock := database.NewRedisClientMock(nil, "", nil)
	corporationService := corpServ.NewService(logger, nil, redisMock)

	a.Error(corporationService.PublishOffers(corporation.NewClientMock([]entity.Offer{}, err), dewoonplaats.Info))
}

func Test_PublishOffers_RedisClientError(t *testing.T) {
	a := assert.New(t)

	err := errors.New("foo")
	logger := logging.NewZapLoggerWithoutSentry()
	redisMock := database.NewRedisClientMock(nil, "", err)
	corporationService := corpServ.NewService(logger, nil, redisMock)

	a.Error(corporationService.PublishOffers(corporation.NewClientMock([]entity.Offer{{}}, nil), dewoonplaats.Info))
}

func Test_PublishOffers_Success_NoOffers(t *testing.T) {
	a := assert.New(t)
	logger := logging.NewZapLoggerWithoutSentry()
	redisMock := database.NewRedisClientMock(nil, "", nil)
	corporationService := corpServ.NewService(logger, nil, redisMock)

	a.Nil(corporationService.PublishOffers(corporation.NewClientMock([]entity.Offer{}, nil), dewoonplaats.Info))
}

func Test_PublishOffers_Success(t *testing.T) {
	a := assert.New(t)
	logger := logging.NewZapLoggerWithoutSentry()
	redisMock := database.NewRedisClientMock(nil, "", nil)
	corporationService := corpServ.NewService(logger, nil, redisMock)

	a.Nil(corporationService.PublishOffers(corporation.NewClientMock([]entity.Offer{{}}, nil), dewoonplaats.Info))
}

func Test_SubscribeOffers_RedisClientError(t *testing.T) {
	a := assert.New(t)
	err := errors.New("foo")
	logger := logging.NewZapLoggerWithoutSentry()
	redisMock := database.NewRedisClientMock(nil, "", err)
	corporationService := corpServ.NewService(logger, nil, redisMock)

	c := make(chan entity.OfferList)
	a.Error(corporationService.SubscribeOffers(c))
}

func Test_SubscribeOffers_Success(t *testing.T) {
	a := assert.New(t)
	logger := logging.NewZapLoggerWithoutSentry()
	corpInfo, err := json.Marshal(entity.OfferList{Corporation: dewoonplaats.Info})
	a.NoError(err)

	redisChan := make(chan *redis.Message)
	go func(corp []byte) {
		redisChan <- &redis.Message{
			Channel: "channel",
			Payload: string(corp),
		}
	}(corpInfo)
	redisMock := database.NewRedisClientMock(redisChan, "", err)
	corporationService := corpServ.NewService(logger, nil, redisMock)

	c := make(chan entity.OfferList)
	go func(c chan entity.OfferList) {
		err := corporationService.SubscribeOffers(c)
		a.NoError(err)
	}(c)

	resultInfo, _ := <-c
	a.Equal(dewoonplaats.Info.Name, resultInfo.Corporation.Name)
	a.Equal(dewoonplaats.Info.Cities, resultInfo.Corporation.Cities)
}
