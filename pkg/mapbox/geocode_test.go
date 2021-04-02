package mapbox_test

import (
	"testing"

	"github.com/woningfinder/woningfinder/internal/bootstrap"

	"github.com/stretchr/testify/assert"
)

func Test_Mapbox_Geocoding_CityDistrictFromAdress(t *testing.T) {
	a := assert.New(t)
	mapboxClient := bootstrap.CreateMapboxClient()
	districtFromAddress, err := mapboxClient.CityDistrictFromAddress("Zuid Esmarkerrondweg 19, Enschede")
	a.NoError(err)
	a.Equal("de leuriks", districtFromAddress)
}

func Test_Mapbox_Geocoding_CityDistrictFromCoords(t *testing.T) {
	a := assert.New(t)
	mapboxClient := bootstrap.CreateMapboxClient()
	districtFromCoords, err := mapboxClient.CityDistrictFromCoords("52.2130417", "6.9075881")
	a.NoError(err)
	a.Equal("hogeland-noord", districtFromCoords)
}
