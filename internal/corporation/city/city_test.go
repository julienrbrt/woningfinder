package city_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/woningfinder/woningfinder/internal/corporation/city"
	"github.com/woningfinder/woningfinder/pkg/logging"
)

func TestCity_Merge(t *testing.T) {
	a := assert.New(t)

	a.Equal((&city.City{Name: "Hengelo OV"}).Merge(), city.Hengelo)
	a.Equal((&city.City{Name: "Hengelo"}).Merge(), city.Hengelo)
	expected := city.City{Name: "a city"}
	a.Equal(expected.Merge(), expected)
}

func TestCity_SuggestedCityDistrictFromName(t *testing.T) {
	a := assert.New(t)

	data := city.SuggestedCityDistrictFromName(logging.NewZapLoggerWithoutSentry(), city.Enschede.Name)
	a.Len(data, len(city.Enschede.Districts()))

	var districts, neighbourhood []string
	for d, n := range data {
		districts = append(districts, d)
		neighbourhood = append(neighbourhood, n...)
	}
	a.ElementsMatch(districts, city.Enschede.Districts())
	a.ElementsMatch(neighbourhood, city.Enschede.Neighbourhoods())

	a.Len(city.SuggestedCityDistrictFromName(logging.NewZapLoggerWithoutSentry(), "unexisting"), 0)
}

func TestCity_HasSuggestedCityDistrict(t *testing.T) {
	a := assert.New(t)
	a.True(city.HasSuggestedCityDistrict(city.Enschede.Name))
	a.False(city.HasSuggestedCityDistrict(city.DeLutte.Name))
	a.False(city.HasSuggestedCityDistrict("unexisting"))
}
