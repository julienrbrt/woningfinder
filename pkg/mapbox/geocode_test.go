package mapbox_test

import (
	"testing"

	"github.com/woningfinder/woningfinder/internal/bootstrap"

	"github.com/stretchr/testify/assert"
)

func Test_Mapbox_Geocoding(t *testing.T) {
	a := assert.New(t)
	mapboxClient := bootstrap.CreateMapboxGeocodingClient()
	districtFromAddress, err := mapboxClient.CityDistrictFromAddress("Hogelandsingel 120 ENSCHEDE")
	districtFromCoords, err := mapboxClient.CityDistrictFromCoords("52.2130417", "6.9075881")
	a.NoError(err)
	a.Equal("hogeland-noord", districtFromAddress)
	a.Equal(districtFromAddress, districtFromCoords)
}
