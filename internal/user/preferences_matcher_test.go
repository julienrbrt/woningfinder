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
			Address:                 "Beatrixstraat 1 R 7161 DJ Neede  A",
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
			Historic:                false,
			CV:                      false,
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
			MinimumPrice:  400,
			MaximumPrice:  950,
			NumberBedroom: 2,
			HasElevator:   true,
		},
	}
}

func Test_MatchPreferences_Location(t *testing.T) {
	a := assert.New(t)
	testUser := getUser()
	testOffer := getOffer()

	a.True(testUser.MatchPreferences(testOffer))

	testUser.HousingPreferences.City = []corporation.City{
		{Name: "Hengelo"},
		{Name: "Enschede",
			District: []corporation.District{
				{Name: "roombeek"},
				{Name: "boddenkamp"},
				{Name: "lasonder-zeggelt"},
			},
		},
	}
	a.False(testUser.MatchPreferences(testOffer))
	testOffer.Housing.City.Name = "Hengelo"
	a.True(testUser.MatchPreferences(testOffer))
	testOffer.Housing.City.Name = "enschede"
	a.False(testUser.MatchPreferences(testOffer))
	testOffer.Housing.City = corporation.City{Name: "Enschede", District: []corporation.District{{Name: "Enschede - Roombeek"}}}
	a.True(testUser.MatchPreferences(testOffer))
}

func Test_MatchPreferences_HousingType(t *testing.T) {
	a := assert.New(t)
	testUser := getUser()
	testOffer := getOffer()

	a.True(testUser.MatchPreferences(testOffer))
	testOffer.Housing.Type = corporation.HousingType{
		Type: corporation.Parking,
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
