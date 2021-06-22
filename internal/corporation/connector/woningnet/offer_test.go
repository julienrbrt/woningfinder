package woningnet_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	bootstrapCorporation "github.com/woningfinder/woningfinder/internal/bootstrap/corporation"
	"github.com/woningfinder/woningfinder/internal/corporation"
	"github.com/woningfinder/woningfinder/internal/corporation/connector/woningnet"
	"github.com/woningfinder/woningfinder/pkg/logging"
	"github.com/woningfinder/woningfinder/pkg/mapbox"
)

func Test_FetchOffer(t *testing.T) {
	a := assert.New(t)
	client := bootstrapCorporation.CreateWoningNetClient(logging.NewZapLoggerWithoutSentry(), mapbox.NewClientMock(nil, "district"), woningnet.HengeloBorneInfo)

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
		a.NotEmpty(offer.Housing.EnergyLabel)
		a.True(offer.Housing.Price > 0)
		a.True(offer.Housing.Size > 0)
		a.True(offer.Housing.NumberBedroom > 0)
		a.True(offer.Housing.BuildingYear > 0)

		a.NotEmpty(offer.SelectionMethod)
		a.NotEmpty(offer.URL)
		a.NotEmpty(offer.ExternalID)
	}
}
