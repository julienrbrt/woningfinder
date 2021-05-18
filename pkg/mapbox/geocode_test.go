package mapbox_test

import (
	"testing"

	"github.com/woningfinder/woningfinder/internal/bootstrap"

	"github.com/stretchr/testify/assert"
)

func Test_Mapbox_Geocoding_CityDistrictFromAdress(t *testing.T) {
	a := assert.New(t)
	mapboxClient := bootstrap.CreateMapboxClient()
	districtFromAddress, err := mapboxClient.CityDistrictFromAddress("Stroinksbleekweg 27 Enschede")
	a.NoError(err)
	a.Equal("roombeek-roomveldje", districtFromAddress)
}
