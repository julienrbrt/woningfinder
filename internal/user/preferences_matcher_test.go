package user_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/woningfinder/woningfinder/internal/corporation"
	"github.com/woningfinder/woningfinder/internal/user"
)

func getOffer() corporation.Offer {
	return corporation.Offer{
		ExternalID: "w758752",
		Housing: corporation.Housing{
			Type: corporation.HousingType{
				Type: corporation.House,
			},
			Latitude:                52.133,
			Longitude:               6.61433,
			Address:                 "Beatrixstraat 1 R 7161 DJ Neede A",
			EnergieLabel:            "A",
			Price:                   656.39,
			Size:                    80,
			NumberRoom:              6,
			NumberBedroom:           2,
			BuildingYear:            2010,
			HousingAllowance:        true,
			Garden:                  false,
			Garage:                  false,
			Elevator:                true,
			Balcony:                 true,
			AccessibilityWheelchair: false,
			AccessibilityScooter:    true,
			Attic:                   false,
		},
		MinAge:        18,
		MaxAge:        35,
		MinIncome:     20000,
		MaxIncome:     28000,
		MinFamilySize: 1,
		MaxFamilySize: 2,
	}
}

func getUser() user.User {
	return user.User{
		BirthYear:    1990,
		YearlyIncome: 30000,
		FamilySize:   3,
		HousingPreferences: user.HousingPreferences{
			Type: []corporation.HousingType{
				corporation.HousingType{
					Type: corporation.House},
				corporation.HousingType{
					Type: corporation.Appartement},
			},
			MaximumPrice:  950,
			NumberBedroom: 2,
			HasElevator:   true,
		},
	}
}

var enschede = corporation.City{Name: "Enschede"}

var hengelo = corporation.City{Name: "Hengelo"}

func Test_MatchPreferences_Location(t *testing.T) {
	a := assert.New(t)
	testUser := getUser()
	testOffer := getOffer()

	a.True(testUser.MatchPreferences(testOffer))
	testUser.HousingPreferences.City = []corporation.City{enschede, hengelo}
	a.False(testUser.MatchPreferences(testOffer))
	testOffer.Housing.City = hengelo
	a.True(testUser.MatchPreferences(testOffer))
	testOffer.Housing.City = enschede
	testUser.HousingPreferences.CityDistrict = []corporation.CityDistrict{
		{
			Name:     "roombeek",
			CityName: "enschede",
		},
		{
			Name:     "boddenkamp",
			CityName: "enschede",
		},
		{
			Name:     "lasonder-zeggelt",
			CityName: "enschede",
		},
	}
	a.False(testUser.MatchPreferences(testOffer))
	testOffer.Housing.CityDistrict = corporation.CityDistrict{CityName: "Enschede", Name: "Enschede - Roombeek"}
	a.True(testUser.MatchPreferences(testOffer))
	testOffer.Housing.CityDistrict = corporation.CityDistrict{CityName: "Enschede", Name: "deppenbroek"}
	a.False(testUser.MatchPreferences(testOffer))
}

func Test_MatchPreferences_HousingType(t *testing.T) {
	a := assert.New(t)
	testUser := getUser()
	testOffer := getOffer()

	a.True(testUser.MatchPreferences(testOffer))
	testOffer.Housing.Type = corporation.HousingType{
		Type: corporation.Undefined,
	}
	a.False(testUser.MatchPreferences(testOffer))
	testUser.HousingPreferences.Type = nil
	a.True(testUser.MatchPreferences(testOffer))
}

func Test_MatchPreferences_Price(t *testing.T) {
	a := assert.New(t)
	testUser := getUser()
	testOffer := getOffer()

	a.True(testUser.MatchPreferences(testOffer))
	testOffer.Housing.Price = 1000
	a.False(testUser.MatchPreferences(testOffer))
}

func Test_MatchPreferences_Criteria(t *testing.T) {
	a := assert.New(t)
	testUser := getUser()
	testOffer := getOffer()

	a.True(testUser.MatchPreferences(testOffer))
	testOffer.Housing.Garden = true
	testUser.HousingPreferences.HasGarden = false
	a.True(testUser.MatchPreferences(testOffer))
	testUser.HousingPreferences.NumberBedroom = 5
	a.False(testUser.MatchPreferences(testOffer))
}
