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
			Type: corporation.House,
			Location: corporation.Location{
				Latitude:  52.133,
				Longitude: 6.61433,
			},
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
		CanApply:   true,
		HasApplied: false,
	}
}

func getUser() user.User {
	return user.User{
		HousingPreferences: user.HousingPreferences{
			Type: []corporation.HousingType{
				corporation.House,
				corporation.Appartement,
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
	testUser.HousingPreferences.District = []string{
		"roombeek",
		"boddenkamp",
		"lasonder-zeggelt",
	}

	testUser.HousingPreferences.City = []string{
		"enschede",
		"hengelo",
	}
	a.False(testUser.MatchPreferences(testOffer))
	testOffer.Housing.District = "Enschede - Roombeek"
	a.False(testUser.MatchPreferences(testOffer))
	testOffer.City.Name = "enschede"
	a.True(testUser.MatchPreferences(testOffer))
}

func Test_MatchPreferences_HousingType(t *testing.T) {
	a := assert.New(t)
	testUser := getUser()
	testOffer := getOffer()

	a.True(testUser.MatchPreferences(testOffer))
	testOffer.Housing.Type = corporation.Parking
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

func Test_MatchPreferences_HasApplied(t *testing.T) {
	a := assert.New(t)
	testUser := getUser()
	testOffer := getOffer()

	a.True(testUser.MatchPreferences(testOffer))
	testOffer.HasApplied = true
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
