package matcher_test

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/woningfinder/woningfinder/internal/corporation"
	"github.com/woningfinder/woningfinder/internal/corporation/dewoonplaats"
	"github.com/woningfinder/woningfinder/internal/database"
	"github.com/woningfinder/woningfinder/internal/domain/entity"
	corporationService "github.com/woningfinder/woningfinder/internal/services/corporation"
	"github.com/woningfinder/woningfinder/internal/services/matcher"
	"github.com/woningfinder/woningfinder/pkg/logging"
)

type mockCorporationService struct {
	corporationService.Service
}

func (m *mockCorporationService) AddCities(cities []entity.City, corporation entity.Corporation) error {
	return nil
}

func Test_PublishOffers_CorporationClientError(t *testing.T) {
	a := assert.New(t)

	err := errors.New("foo")
	logger := logging.NewZapLoggerWithoutSentry()
	redisMock := database.NewRedisClientMock("", nil, nil)
	matcherService := matcher.NewService(logger, redisMock, nil, nil, nil)

	a.Error(matcherService.PublishOffers(corporation.NewClientMock([]entity.Offer{}, err), dewoonplaats.Info))
}

func Test_PublishOffers_RedisClientError(t *testing.T) {
	a := assert.New(t)

	err := errors.New("foo")
	logger := logging.NewZapLoggerWithoutSentry()
	redisMock := database.NewRedisClientMock("", nil, err)
	matcherService := matcher.NewService(logger, redisMock, nil, nil, nil)

	a.Error(matcherService.PublishOffers(corporation.NewClientMock([]entity.Offer{{}}, nil), dewoonplaats.Info))
}

func Test_PublishOffers_Success_NoOffers(t *testing.T) {
	a := assert.New(t)
	logger := logging.NewZapLoggerWithoutSentry()
	redisMock := database.NewRedisClientMock("", nil, nil)
	matcherService := matcher.NewService(logger, redisMock, nil, nil, nil)

	a.Nil(matcherService.PublishOffers(corporation.NewClientMock([]entity.Offer{}, nil), dewoonplaats.Info))
}

func Test_PublishOffers_Success(t *testing.T) {
	a := assert.New(t)
	logger := logging.NewZapLoggerWithoutSentry()
	redisMock := database.NewRedisClientMock("", nil, nil)
	matcherService := matcher.NewService(logger, redisMock, nil, &mockCorporationService{}, nil)

	a.Nil(matcherService.PublishOffers(corporation.NewClientMock([]entity.Offer{{}}, nil), dewoonplaats.Info))
}

func Test_SubscribeOffers_RedisClientError(t *testing.T) {
	a := assert.New(t)
	err := errors.New("foo")
	logger := logging.NewZapLoggerWithoutSentry()
	redisMock := database.NewRedisClientMock("", nil, err)
	matcherService := matcher.NewService(logger, redisMock, nil, nil, nil)

	c := make(chan entity.OfferList)
	a.Error(matcherService.SubscribeOffers(c))
}

func Test_SubscribeOffers_Success(t *testing.T) {
	a := assert.New(t)
	logger := logging.NewZapLoggerWithoutSentry()
	corpInfo, err := json.Marshal(entity.OfferList{Corporation: dewoonplaats.Info})
	a.NoError(err)

	redisMock := database.NewRedisClientMock("", []string{string(corpInfo)}, err)
	matcherService := matcher.NewService(logger, redisMock, nil, nil, nil)

	c := make(chan entity.OfferList)
	go func(c chan entity.OfferList) {
		err := matcherService.SubscribeOffers(c)
		a.NoError(err)
	}(c)

	resultInfo, _ := <-c
	a.Equal(dewoonplaats.Info.Name, resultInfo.Corporation.Name)
	a.Equal(dewoonplaats.Info.Cities, resultInfo.Corporation.Cities)
}
