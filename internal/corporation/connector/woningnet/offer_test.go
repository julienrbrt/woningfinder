package woningnet_test

import (
	"testing"

	bootstrapCorporation "github.com/julienrbrt/woningfinder/internal/bootstrap/corporation"
	"github.com/julienrbrt/woningfinder/internal/corporation"
	"github.com/julienrbrt/woningfinder/internal/corporation/connector/woningnet"
	"github.com/julienrbrt/woningfinder/pkg/logging"
	"github.com/julienrbrt/woningfinder/pkg/mapbox"
	"github.com/stretchr/testify/assert"
)

func Test_FetchOffers(t *testing.T) {
	a := assert.New(t)

	corporations := []corporation.Corporation{
		woningnet.UtrechtInfo,
		woningnet.AmsterdamInfo,
		woningnet.AlmereInfo,
		woningnet.WoonkeusInfo,
		woningnet.EemvalleiInfo,
		woningnet.WoonserviceInfo,
		woningnet.MiddenHollandInfo,
		woningnet.BovenGroningenInfo,
		woningnet.GooiVechtstreekInfo,
		woningnet.GroningenInfo,
		woningnet.HuiswaartsInfo,
		woningnet.WoongaardInfo,
	}

	for _, c := range corporations {
		// note for testing mapbox city parsing please use the bootstrap client instead of the mock
		client := bootstrapCorporation.CreateWoningNetClient(logging.NewZapLoggerWithoutSentry(), mapbox.NewClientMock(nil, "district"), c)

		ch := make(chan corporation.Offer, 1000)
		a.NoError(client.FetchOffers(ch))
		for offer := range ch {
			a.NotEmpty(offer.CorporationName)

			// verify housing validity
			a.NotEmpty(offer.Housing.Type)
			if offer.Housing.Type == corporation.HousingTypeUndefined {
				a.FailNow("should not be undefined")
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

			// if passes stop test for a corporation
			if !t.Failed() {
				break
			}
		}
	}
}
