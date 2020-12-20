package itris_test

import (
	"testing"

	"github.com/woningfinder/woningfinder/internal/corporation"

	"github.com/woningfinder/woningfinder/internal/bootstrap"

	"github.com/stretchr/testify/assert"
)

func Test_FetchOffer(t *testing.T) {
	a := assert.New(t)
	itrisConnector := bootstrap.CreateOnsHuisClient()

	offers, err := itrisConnector.FetchOffer()
	a.NoError(err)
	a.True(len(offers) > 0)
	for _, offer := range offers {
		a.NotEmpty(offer.Housing.Type.Type)
		if offer.Housing.Type.Type == corporation.Undefined {
			continue
		}

		a.NotNil(offer.SelectionMethod)
		a.NotNil(offer.SelectionDate)
		a.NotEmpty(offer.URL)
		a.NotEmpty(offer.ExternalID)

		a.NotNil(offer.Housing)
		a.NotEmpty(offer.Housing.City.Name)
		a.NotEmpty(offer.Housing.Address)
		a.NotEmpty(offer.Housing.EnergieLabel)
		a.True(offer.Housing.BuildingYear > 0)
		if offer.Housing.Size == 0 {
			continue
		}
		a.True(offer.Housing.Size > 0)
		a.True(offer.Housing.Price > 0)
		a.True(offer.Housing.Longitude > 0)
		a.True(offer.Housing.Latitude > 0)
		a.True(offer.Housing.NumberRoom > 0)

		// test only for one offer
		return
	}
}
