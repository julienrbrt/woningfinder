package mapbox_test

import (
	"testing"

	"github.com/woningfinder/woningfinder/internal/bootstrap"

	"github.com/stretchr/testify/assert"
)

func Test_Mapbox_Geocoding_CityDistrictFromAdress(t *testing.T) {
	a := assert.New(t)
	mapboxClient := bootstrap.CreateMapboxClient()
	districtFromAddress, err := mapboxClient.CityDistrictFromAddress("Oogstplein 26, 7545 HP Enschede")
	a.NoError(err)
	a.Equal("stevenfenne", districtFromAddress)
}
