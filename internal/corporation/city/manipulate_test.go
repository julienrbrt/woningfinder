package city_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/woningfinder/woningfinder/internal/corporation"
	"github.com/woningfinder/woningfinder/internal/corporation/city"
	"github.com/woningfinder/woningfinder/pkg/logging"
)

func TestCity_Merge(t *testing.T) {
	a := assert.New(t)

	mergedCity := city.Merge(corporation.City{Name: "Hengelo OV"})
	a.Equal(mergedCity, city.Hengelo)

	mergedCity = city.Merge(corporation.City{Name: "Hengelo"})
	a.Equal(mergedCity, city.Hengelo)

	expected := corporation.City{Name: "a city"}
	mergedCity = city.Merge(expected)
	a.Equal(mergedCity, expected)
}

func TestCity_SuggestedCityDistrictFromName(t *testing.T) {
	a := assert.New(t)

	districts := city.SuggestedCityDistrictFromName(logging.NewZapLoggerWithoutSentry(), city.Enschede.Name)
	a.Equal(districts, city.Enschede.District)

	districts = city.SuggestedCityDistrictFromName(logging.NewZapLoggerWithoutSentry(), "unexisting")
	a.Len(districts, 0)
}

func TestCity_HasSuggestedCityDistrict(t *testing.T) {
	a := assert.New(t)

	ok := city.HasSuggestedCityDistrict(city.Enschede.Name)
	a.True(ok)

	ok = city.HasSuggestedCityDistrict(city.DeLutte.Name)
	a.False(ok)

	ok = city.HasSuggestedCityDistrict("unexisting")
	a.False(ok)
}
