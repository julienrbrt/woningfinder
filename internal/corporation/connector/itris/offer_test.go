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

func Test_FetchOffer(t *testing.T) {
	a := assert.New(t)
	mockMapbox := mapbox.NewClientMock(nil, "district")

	for _, clientURL := range clientList {
		// note for testing mapbox city parsing please use the bootstrap client instead of the mock
		client, err := itris.NewClient(logging.NewZapLoggerWithoutSentry(), mockMapbox, clientURL)
		a.NoError(err)

		offers, err := client.GetOffers()
		a.NoError(err)
		a.True(len(offers) > 0)

		for _, offer := range offers {
			// verify housing validity
			a.NotEmpty(offer.Housing.Type)
			if offer.Housing.Type == corporation.HousingTypeUndefined {
				continue
			}

			a.NotEmpty(offer.Housing.Address)
			a.NotEmpty(offer.Housing.CityName)
			a.NotEmpty(offer.Housing.CityDistrict)
			a.True(offer.Housing.Price > 0)
			a.True(offer.Housing.Size > 0)
			a.True(offer.Housing.NumberBedroom > 0)
			a.True(offer.Housing.BuildingYear > 0)

			a.NotEmpty(offer.SelectionMethod)
			a.NotEmpty(offer.URL)
			a.NotEmpty(offer.RawPictureURL)
			a.NotEmpty(offer.ExternalID)
		}
	}
}
