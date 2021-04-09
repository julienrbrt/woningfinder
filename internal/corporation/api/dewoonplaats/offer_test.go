package dewoonplaats_test

import (
	"testing"

	"github.com/woningfinder/woningfinder/pkg/mapbox"

	"github.com/woningfinder/woningfinder/pkg/logging"

	"github.com/woningfinder/woningfinder/internal/bootstrap"
	"github.com/woningfinder/woningfinder/internal/entity"

	"github.com/stretchr/testify/assert"
)

func Test_FetchOffer(t *testing.T) {
	a := assert.New(t)
	client := bootstrap.CreateDeWoonplaatsClient(logging.NewZapLoggerWithoutSentry(), mapbox.NewClientMock(nil, "district"))

	offers, err := client.FetchOffer()
	a.NoError(err)
	a.True(len(offers) > 0)
	for _, offer := range offers {
		// verify housing validity
		a.NotEmpty(offer.Housing.Type)
		if offer.Housing.Type == entity.HousingTypeUndefined {
			continue
		}

		a.NotEmpty(offer.Housing.Address)
		a.NotEmpty(offer.Housing.City.Name)
		a.NotEmpty(offer.Housing.CityDistrict)
		a.NotEmpty(offer.Housing.EnergieLabel)
		a.True(offer.Housing.Price > 0)
		// a.True(offer.Housing.Size > 0)
		a.True(offer.Housing.NumberRoom > 0)
		a.True(offer.Housing.NumberBedroom > 0)
		a.True(offer.Housing.BuildingYear > 0)

		a.NotEmpty(offer.SelectionMethod)
		a.NotNil(offer.SelectionDate)
		a.NotEmpty(offer.URL)
		a.NotEmpty(offer.ExternalID)

		// test only for one offer
		return
	}
}
