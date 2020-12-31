package itris_test

import (
	"testing"

	"github.com/woningfinder/woningfinder/internal/bootstrap"

	"github.com/woningfinder/woningfinder/pkg/logging"

	"github.com/woningfinder/woningfinder/internal/corporation"
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

	for _, url := range clientList {
		client := itris.NewClient(logging.NewZapLogger(), bootstrap.CreateMapboxGeocodingClient(), url)

		offers, err := client.FetchOffer()
		a.NoError(err)
		if !hadOffer {
			hadOffer = len(offers) > 0
		}
		for _, offer := range offers {
			a.True(offer.Housing.IsValid())

			if offer.Housing.Type.Type == corporation.Undefined {
				continue
			}

			a.NotNil(offer.SelectionMethod)
			a.NotNil(offer.SelectionDate)
			a.NotEmpty(offer.URL)
			a.NotEmpty(offer.ExternalID)

			// test only for one offer
			break
		}
	}

	a.True(hadOffer)
}
