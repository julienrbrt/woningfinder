package matcher_test

import (
	"encoding/json"
	"errors"
	"testing"

	bootstrapCorporation "github.com/julienrbrt/woningfinder/internal/bootstrap/corporation"
	"github.com/julienrbrt/woningfinder/internal/corporation"
	"github.com/julienrbrt/woningfinder/internal/corporation/city"
	"github.com/julienrbrt/woningfinder/internal/corporation/connector"
	"github.com/julienrbrt/woningfinder/internal/customer/matcher"
	"github.com/julienrbrt/woningfinder/internal/database"
	corporationService "github.com/julienrbrt/woningfinder/internal/services/corporation"
	matcherService "github.com/julienrbrt/woningfinder/internal/services/matcher"
	"github.com/julienrbrt/woningfinder/pkg/logging"
	"github.com/stretchr/testify/assert"
)

type mockCorporationService struct {
	corporationService.Service
}

func (m *mockCorporationService) LinkCities(cities []city.City, hasLocation bool, corporation ...corporation.Corporation) error {
	return nil
}

var corporationInfo = corporation.Corporation{
	Name: "OnsHuis",
	URL:  "https://example.com",
	Cities: []city.City{
		city.Enschede,
		city.Hengelo,
	},
	SelectionMethod: []corporation.SelectionMethod{
		corporation.SelectionRandom,
	},
}

var (
	cities = bootstrapCorporation.CreateConnectorProvider(logging.NewZapLoggerWithoutSentry(), nil).GetCities()
	m      = matcher.NewMatcher(city.NewSuggester(cities))
)

func Test_SendOffers_RedisClientError(t *testing.T) {
	a := assert.New(t)

	err := errors.New("foo")
	logger := logging.NewZapLoggerWithoutSentry()
	redisMock := database.NewRedisClientMock("", nil, err)
	matcherService := matcherService.NewService(logger, redisMock, nil, nil, nil, nil, m, nil)

	a.Error(matcherService.SendOffers(corporation.Offers{CorporationName: corporationInfo.Name}))
}

func Test_SendOffers_Success_EmptyOffer(t *testing.T) {
	a := assert.New(t)
	logger := logging.NewZapLoggerWithoutSentry()
	redisMock := database.NewRedisClientMock("", nil, nil)
	matcherService := matcherService.NewService(logger, redisMock, nil, nil, nil, nil, m, connector.NewConnectorProvider([]connector.Provider{{Corporation: corporationInfo}}))

	a.Nil(matcherService.SendOffers(corporation.Offers{}))
}

func Test_SendOffers_Success(t *testing.T) {
	a := assert.New(t)
	logger := logging.NewZapLoggerWithoutSentry()
	redisMock := database.NewRedisClientMock("", nil, nil)
	matcherService := matcherService.NewService(logger, redisMock, nil, nil, &mockCorporationService{}, nil, m, connector.NewConnectorProvider([]connector.Provider{{Corporation: corporationInfo}}))

	a.Nil(matcherService.SendOffers(corporation.Offers{CorporationName: corporationInfo.Name}))
}

func Test_RetrieveOffers_RedisClientError(t *testing.T) {
	a := assert.New(t)
	err := errors.New("foo")
	logger := logging.NewZapLoggerWithoutSentry()
	redisMock := database.NewRedisClientMock("", nil, err)
	matcherService := matcherService.NewService(logger, redisMock, nil, nil, nil, nil, m, nil)

	c := make(chan corporation.Offers)
	a.Error(matcherService.RetrieveOffers(c))
}

func Test_RetrieveOffers_Success(t *testing.T) {
	a := assert.New(t)
	logger := logging.NewZapLoggerWithoutSentry()
	corpInfo, err := json.Marshal(corporation.Offer{CorporationName: corporationInfo.Name})
	a.NoError(err)

	redisMock := database.NewRedisClientMock("", []string{string(corpInfo)}, err)
	matcherService := matcherService.NewService(logger, redisMock, nil, nil, nil, nil, m, nil)

	c := make(chan corporation.Offers)
	go func(c chan corporation.Offers) {
		err := matcherService.RetrieveOffers(c)
		a.NoError(err)
	}(c)

	resultInfo := <-c
	a.Equal(corporationInfo.Name, resultInfo.CorporationName)
}
