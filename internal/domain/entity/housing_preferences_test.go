package entity_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/woningfinder/woningfinder/internal/domain/entity"
)

func Test_HousingPreferences_IsValid(t *testing.T) {
	a := assert.New(t)

	housingPreferences := entity.HousingPreferences{
		Type: []entity.HousingType{
			{
				Type: entity.HousingTypeAppartement,
			},
		},
		City: []entity.City{
			{
				Name: "Enschede",
			},
		},
	}

	a.Nil(housingPreferences.IsValid())
}

func Test_HousingPreferences_IsValid_Invalid(t *testing.T) {
	a := assert.New(t)
	housingPreferences := entity.HousingPreferences{
		City: []entity.City{
			{
				Name: "Enschede",
			},
		},
	}
	a.Error(housingPreferences.IsValid())
	housingPreferences.City = nil
	a.Error(housingPreferences.IsValid())
}
