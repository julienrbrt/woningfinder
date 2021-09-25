package city_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/woningfinder/woningfinder/internal/city"
)

func TestManipulate_MergeCity(t *testing.T) {
	a := assert.New(t)

	a.Equal((city.MergeCity(city.City{Name: "Hengelo OV"})), city.Hengelo)
	a.Equal(city.MergeCity(city.City{Name: "Hengelo"}), city.Hengelo)
	expected := city.City{Name: "a city"}
	a.Equal(city.MergeCity(expected), expected)
}

func TestManipulate_SuggestedCityDistrictFromName(t *testing.T) {
	a := assert.New(t)

	data, ok := city.SuggestedCityDistrictFromName(city.Enschede.Name)
	a.True(ok)
	a.Len(data, len(city.Enschede.Districts()))

	var districts, neighbourhood []string
	for d, n := range data {
		districts = append(districts, d)
		neighbourhood = append(neighbourhood, n...)
	}
	a.ElementsMatch(districts, city.Enschede.Districts())
	a.ElementsMatch(neighbourhood, city.Enschede.Neighbourhoods())

	_, ok = city.SuggestedCityDistrictFromName("unexisting")
	a.False(ok)
}

func TestManipulate_HasSuggestedCityDistrict(t *testing.T) {
	a := assert.New(t)
	a.True(city.HasSuggestedCityDistrict(city.Enschede.Name))
	a.False(city.HasSuggestedCityDistrict(city.DeLutte.Name))
	a.False(city.HasSuggestedCityDistrict("unexisting"))
}
