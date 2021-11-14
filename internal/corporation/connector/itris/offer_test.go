package itris_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/woningfinder/woningfinder/internal/corporation"
	"github.com/woningfinder/woningfinder/internal/corporation/connector/itris"
	"github.com/woningfinder/woningfinder/pkg/logging"
	"github.com/woningfinder/woningfinder/pkg/mapbox"
)

var clientList = []corporation.Corporation{
	itris.OnsHuisInfo,
	itris.MijandeWonenInfo,
}

func Test_FetchOffers(t *testing.T) {
	a := assert.New(t)
	mockMapbox := mapbox.NewClientMock(nil, "district")

	for _, clientURL := range clientList {
		// note for testing mapbox city parsing please use the bootstrap client instead of the mock
		client, err := itris.NewClient(logging.NewZapLoggerWithoutSentry(), mockMapbox, clientURL)
		a.NoError(err)

		ch := make(chan corporation.Offer, 1000)
		a.NoError(client.FetchOffers(ch))
		for offer := range ch {
			a.NotEmpty(offer.CorporationName)

			// verify housing validity
			a.NotEmpty(offer.Housing.Type)
			if offer.Housing.Type == corporation.HousingTypeUndefined {
				continue
			}

			a.NotEmpty(offer.Housing.Address)
			a.NotEmpty(offer.Housing.CityName)
			a.NotEmpty(offer.Housing.CityDistrict)
			a.True(offer.Housing.Price > 0)
			a.True(offer.Housing.Size >= 0)
			a.True(offer.Housing.NumberBedroom > 0)

			a.NotEmpty(offer.URL)
			a.NotEmpty(offer.RawPictureURL)
			a.NotEmpty(offer.ExternalID)

			// if passes stop test for the corporation
			if !t.Failed() {
				break
			}
		}
	}
}
