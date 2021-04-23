package onshuis_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/woningfinder/woningfinder/internal/corporation"
	"github.com/woningfinder/woningfinder/internal/corporation/connector/itris"
	"github.com/woningfinder/woningfinder/internal/corporation/connector/itris/onshuis"
	"github.com/woningfinder/woningfinder/pkg/logging"
	"github.com/woningfinder/woningfinder/pkg/mapbox"
)

func Test_FetchOffer(t *testing.T) {
	a := assert.New(t)
	mockMapbox := mapbox.NewClientMock(nil, "district")

	client, err := itris.NewClient(logging.NewZapLoggerWithoutSentry(), mockMapbox, "https://mijn.onshuis.com", onshuis.DetailsParser)
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
		a.NotEmpty(offer.Housing.City.Name)
		a.NotEmpty(offer.Housing.CityDistrict)
		a.NotEmpty(offer.Housing.EnergieLabel)
		a.True(offer.Housing.Price > 0)
		a.True(offer.Housing.Size > 0)
		a.True(offer.Housing.NumberRoom > 0)
		a.True(offer.Housing.NumberBedroom > 0) // might fail randomly if there is only studios available
		a.True(offer.Housing.BuildingYear > 0)

		a.NotEmpty(offer.SelectionMethod)
		a.NotNil(offer.SelectionDate)
		a.NotEmpty(offer.URL)
		a.NotEmpty(offer.ExternalID)

		// test only for one offer
		break
	}
}
