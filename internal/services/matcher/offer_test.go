package matcher_test

import (
	"encoding/json"
	"errors"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/woningfinder/woningfinder/internal/corporation"
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

var corporationInfo = entity.Corporation{
	APIEndpoint: &url.URL{Scheme: "https", Host: "example.com"},
	Name:        "OnsHuis",
	URL:         "https://example.com",
	Cities: []entity.City{
		{Name: "Enschede"},
		{Name: "Hengelo"},
	},
	SelectionMethod: []entity.SelectionMethod{
		entity.SelectionRandom,
	},
}

func Test_PublishOffers_CorporationClientError(t *testing.T) {
	a := assert.New(t)

	err := errors.New("foo")
	logger := logging.NewZapLoggerWithoutSentry()
	redisMock := database.NewRedisClientMock("", nil, nil)
	matcherService := matcher.NewService(logger, redisMock, nil, nil, nil)

	a.Error(matcherService.PublishOffers(corporation.NewClientMock([]entity.Offer{}, err), corporationInfo))
}

func Test_PublishOffers_RedisClientError(t *testing.T) {
	a := assert.New(t)

	err := errors.New("foo")
	logger := logging.NewZapLoggerWithoutSentry()
	redisMock := database.NewRedisClientMock("", nil, err)
	matcherService := matcher.NewService(logger, redisMock, nil, nil, nil)

	a.Error(matcherService.PublishOffers(corporation.NewClientMock([]entity.Offer{{}}, nil), corporationInfo))
}

func Test_PublishOffers_Success_NoOffers(t *testing.T) {
	a := assert.New(t)
	logger := logging.NewZapLoggerWithoutSentry()
	redisMock := database.NewRedisClientMock("", nil, nil)
	matcherService := matcher.NewService(logger, redisMock, nil, nil, nil)

	a.Nil(matcherService.PublishOffers(corporation.NewClientMock([]entity.Offer{}, nil), corporationInfo))
}

func Test_PublishOffers_Success(t *testing.T) {
	a := assert.New(t)
	logger := logging.NewZapLoggerWithoutSentry()
	redisMock := database.NewRedisClientMock("", nil, nil)
	matcherService := matcher.NewService(logger, redisMock, nil, &mockCorporationService{}, nil)

	a.Nil(matcherService.PublishOffers(corporation.NewClientMock([]entity.Offer{{}}, nil), corporationInfo))
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
	corpInfo, err := json.Marshal(entity.OfferList{Corporation: corporationInfo})
	a.NoError(err)

	redisMock := database.NewRedisClientMock("", []string{string(corpInfo)}, err)
	matcherService := matcher.NewService(logger, redisMock, nil, nil, nil)

	c := make(chan entity.OfferList)
	go func(c chan entity.OfferList) {
		err := matcherService.SubscribeOffers(c)
		a.NoError(err)
	}(c)

	resultInfo := <-c
	a.Equal(corporationInfo.Name, resultInfo.Corporation.Name)
	a.Equal(corporationInfo.Cities, resultInfo.Corporation.Cities)
}
