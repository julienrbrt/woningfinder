package itris_test

import (
	"testing"

	"github.com/woningfinder/woningfinder/pkg/mapbox"

	"github.com/woningfinder/woningfinder/internal/entity"

	"github.com/woningfinder/woningfinder/pkg/logging"

	"github.com/woningfinder/woningfinder/internal/corporation/connector/itris"

	"github.com/stretchr/testify/assert"
)

var clientList = []string{
	"https://mijn.onshuis.com",
	"https://mijn.mijande.nl",
	"https://mijn.woongoedzvl.nl",
}

func Test_FetchOffer(t *testing.T) {
	a := assert.New(t)
	// checks if at least one of the tested housing corporation had offers
	hadOffer := false
	mockMapbox := mapbox.NewClientMock(nil, "district")

	for _, url := range clientList {
		client, err := itris.NewClient(logging.NewZapLoggerWithoutSentry(), mockMapbox, url)
		a.NoError(err)

		offers, err := client.FetchOffer()
		a.NoError(err)
		if !hadOffer {
			hadOffer = len(offers) > 0
		}
		for _, offer := range offers {
			// verify housing validity
			a.NotEmpty(offer.Housing.Type)
			if offer.Housing.Type == entity.HousingTypeUndefined {
				continue
			}

			a.NotEmpty(offer.Housing.Address)
			a.NotEmpty(offer.Housing.City.Name)
			// a.NotEmpty(offer.Housing.CityDistrict)
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
			break
		}
	}

	a.True(hadOffer)
}
