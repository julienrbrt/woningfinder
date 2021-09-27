package matcher_test

import (
	"encoding/json"
	"errors"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/woningfinder/woningfinder/internal/city"
	"github.com/woningfinder/woningfinder/internal/corporation"
	"github.com/woningfinder/woningfinder/internal/corporation/connector"
	"github.com/woningfinder/woningfinder/internal/customer/matcher"
	"github.com/woningfinder/woningfinder/internal/database"
	corporationService "github.com/woningfinder/woningfinder/internal/services/corporation"
	matcherService "github.com/woningfinder/woningfinder/internal/services/matcher"
	"github.com/woningfinder/woningfinder/pkg/logging"
)

type mockCorporationService struct {
	corporationService.Service
}

func (m *mockCorporationService) LinkCities(cities []city.City, corporation ...corporation.Corporation) error {
	return nil
}

var corporationInfo = corporation.Corporation{
	APIEndpoint: &url.URL{Scheme: "https", Host: "example.com"},
	Name:        "OnsHuis",
	URL:         "https://example.com",
	Cities: []city.City{
		city.Enschede,
		city.Hengelo,
	},
	SelectionMethod: []corporation.SelectionMethod{
		corporation.SelectionRandom,
	},
}

func Test_PushOffers_CorporationClientError(t *testing.T) {
	a := assert.New(t)

	err := errors.New("foo")
	logger := logging.NewZapLoggerWithoutSentry()
	redisMock := database.NewRedisClientMock("", nil, nil)
	matcherService := matcherService.NewService(logger, redisMock, nil, nil, nil, nil, matcher.NewMatcher(), nil)

	a.Error(matcherService.PushOffers(connector.NewClientMock([]corporation.Offer{}, err), corporationInfo))
}

func Test_PushOffers_RedisClientError(t *testing.T) {
	a := assert.New(t)

	err := errors.New("foo")
	logger := logging.NewZapLoggerWithoutSentry()
	redisMock := database.NewRedisClientMock("", nil, err)
	matcherService := matcherService.NewService(logger, redisMock, nil, nil, nil, nil, matcher.NewMatcher(), nil)

	a.Error(matcherService.PushOffers(connector.NewClientMock([]corporation.Offer{{}}, nil), corporationInfo))
}

func Test_PushOffers_Success_NoOffers(t *testing.T) {
	a := assert.New(t)
	logger := logging.NewZapLoggerWithoutSentry()
	redisMock := database.NewRedisClientMock("", nil, nil)
	matcherService := matcherService.NewService(logger, redisMock, nil, nil, nil, nil, matcher.NewMatcher(), nil)

	a.Nil(matcherService.PushOffers(connector.NewClientMock([]corporation.Offer{}, nil), corporationInfo))
}

func Test_PushOffers_Success(t *testing.T) {
	a := assert.New(t)
	logger := logging.NewZapLoggerWithoutSentry()
	redisMock := database.NewRedisClientMock("", nil, nil)
	matcherService := matcherService.NewService(logger, redisMock, nil, nil, &mockCorporationService{}, nil, matcher.NewMatcher(), nil)

	a.Nil(matcherService.PushOffers(connector.NewClientMock([]corporation.Offer{{}}, nil), corporationInfo))
}

func Test_SubscribeOffers_RedisClientError(t *testing.T) {
	a := assert.New(t)
	err := errors.New("foo")
	logger := logging.NewZapLoggerWithoutSentry()
	redisMock := database.NewRedisClientMock("", nil, err)
	matcherService := matcherService.NewService(logger, redisMock, nil, nil, nil, nil, matcher.NewMatcher(), nil)

	c := make(chan corporation.Offers)
	a.Error(matcherService.SubscribeOffers(c))
}

func Test_SubscribeOffers_Success(t *testing.T) {
	a := assert.New(t)
	logger := logging.NewZapLoggerWithoutSentry()
	corpInfo, err := json.Marshal(corporation.Offers{Corporation: corporationInfo})
	a.NoError(err)

	redisMock := database.NewRedisClientMock("", []string{string(corpInfo)}, err)
	matcherService := matcherService.NewService(logger, redisMock, nil, nil, nil, nil, matcher.NewMatcher(), nil)

	c := make(chan corporation.Offers)
	go func(c chan corporation.Offers) {
		err := matcherService.SubscribeOffers(c)
		a.NoError(err)
	}(c)

	resultInfo := <-c
	a.Equal(corporationInfo.Name, resultInfo.Corporation.Name)
	a.Equal(corporationInfo.Cities, resultInfo.Corporation.Cities)
}
