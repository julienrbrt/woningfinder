package dewoonplaats_test

import (
	"testing"

	"github.com/woningfinder/woningfinder/pkg/logging"

	"github.com/woningfinder/woningfinder/internal/corporation"

	"github.com/woningfinder/woningfinder/internal/bootstrap"

	"github.com/stretchr/testify/assert"
)

func Test_FetchOffer(t *testing.T) {
	a := assert.New(t)
	client := bootstrap.CreateDeWoonplaatsClient(logging.NewZapLogger())

	offers, err := client.FetchOffer()
	a.NoError(err)
	a.True(len(offers) > 0)
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
		return
	}
}
