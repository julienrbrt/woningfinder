package corporation_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/woningfinder/woningfinder/internal/corporation"
)

func Test_Housing_IsValid_InvalidHousing(t *testing.T) {
	a := assert.New(t)
	housing := corporation.Housing{
		Type: corporation.HousingType{
			Type: corporation.House,
		},
		Latitude:  52.23148,
		Longitude: 6.89277,
	}
	a.False(housing.IsValid())
}

func Test_Housing_IsValid_ValidUndefined(t *testing.T) {
	a := assert.New(t)
	housing := corporation.Housing{
		Type: corporation.HousingType{
			Type: corporation.Undefined,
		},
		Latitude:  52.23148,
		Longitude: 6.89277,
	}
	a.True(housing.IsValid())
}

func Test_Housing_IsValid_Valid(t *testing.T) {
	a := assert.New(t)
	housing := corporation.Housing{
		Type: corporation.HousingType{
			Type: corporation.House,
		},
		City: corporation.City{
			Name: "enschede",
		},
		CityDistrict: corporation.CityDistrict{
			CityName: "enschede",
			Name:     "roombeek",
		},
		Address:          "Beatrixstraat 1 R 7161 DJ Neede A",
		EnergieLabel:     "A",
		Price:            656.39,
		Size:             80,
		NumberRoom:       6,
		NumberBedroom:    2,
		BuildingYear:     2010,
		HousingAllowance: true,
		Garden:           false,
		Garage:           false,
		Elevator:         true,
		Balcony:          true,
		Attic:            false,
	}
	a.True(housing.IsValid())
}
