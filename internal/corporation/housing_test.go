package corporation_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/woningfinder/woningfinder/internal/corporation"
)

func Test_Offer_SetDistrict(t *testing.T) {
	a := assert.New(t)
	housing := corporation.Housing{
		Type: corporation.HousingType{
			Type: corporation.House,
		},
		Latitude:  52.23148,
		Longitude: 6.89277,
	}
	housing.SetDistrict()
	a.Len(housing.City.District, 1)
	a.Equal("roombeek", housing.City.District[0].Name)
}

func Test_Offer_SetDistrict_Empty(t *testing.T) {
	a := assert.New(t)
	housing := corporation.Housing{
		Type: corporation.HousingType{
			Type: corporation.House,
		},
		Latitude: 0,
	}
	housing.SetDistrict()
	a.Nil(housing.City.District)
}

func Test_Offer_SetDistrict_AlreadySet(t *testing.T) {
	a := assert.New(t)
	housing := corporation.Housing{
		Type: corporation.HousingType{
			Type: corporation.House,
		},
		City: corporation.City{
			District: []corporation.District{
				{Name: "set"},
			},
		},
	}
	housing.SetDistrict()
	a.Equal(housing.City.District[0].Name, "set")
}
