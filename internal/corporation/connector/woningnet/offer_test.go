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

	corporations := []corporation.Corporation{
		woningnet.UtrechtInfo,
		woningnet.HengeloBorneInfo,
		woningnet.AmsterdamInfo,
		woningnet.AlmereInfo,
		woningnet.WoonkeusInfo,
		woningnet.EemvalleiInfo,
		woningnet.WoonserviceInfo,
		woningnet.MercatusInfo,
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
