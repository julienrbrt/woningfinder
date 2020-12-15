package corporation_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/woningfinder/woningfinder/internal/corporation"
)

func Test_Offer_SetCityDistrict(t *testing.T) {
	a := assert.New(t)
	housing := corporation.Housing{
		Type: corporation.HousingType{
			Type: corporation.House,
		},
		Latitude:  52.23148,
		Longitude: 6.89277,
	}
	housing.SetCityDistrict()
	a.Equal("roombeek", housing.CityDistrict.Name)
}

func Test_Offer_SetCityDistrict_Empty(t *testing.T) {
	a := assert.New(t)
	housing := corporation.Housing{
		Type: corporation.HousingType{
			Type: corporation.House,
		},
		Latitude: 0,
	}
	housing.SetCityDistrict()
	a.Equal(housing.CityDistrict.Name, "")
}

func Test_Offer_SetCityDistrict_AlreadySet(t *testing.T) {
	a := assert.New(t)
	housing := corporation.Housing{
		Type: corporation.HousingType{
			Type: corporation.House,
		},
		CityDistrict: corporation.CityDistrict{
			Name: "set",
		},
	}
	housing.SetCityDistrict()
	a.Equal("set", housing.CityDistrict.Name)
}
