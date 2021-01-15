package entity_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/woningfinder/woningfinder/internal/domain/entity"
)

func getOffer() entity.Offer {
	return entity.Offer{
		ExternalID: "w758752",
		Housing: entity.Housing{
			Type: entity.HousingType{
				Type: entity.HousingTypeHouse,
			},
			Latitude:         52.133,
			Longitude:        6.61433,
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
			Accessible:       true,
		},
		MinAge:        18,
		MaxAge:        35,
		MinIncome:     20000,
		MaxIncome:     28000,
		MinFamilySize: 1,
		MaxFamilySize: 2,
	}
}

func getUser() entity.User {
	return entity.User{
		Name:         "Test",
		Email:        "test@example.org",
		BirthYear:    1990,
		YearlyIncome: 30000,
		FamilySize:   3,
		Tier: entity.Tier{
			Name: entity.PlanZeker,
		},
		HousingPreferences: []entity.HousingPreferences{
			{
				Type: []entity.HousingType{
					entity.HousingType{
						Type: entity.HousingTypeHouse},
					entity.HousingType{
						Type: entity.HousingTypeAppartement},
				},
				MaximumPrice:  950,
				NumberBedroom: 2,
				HasElevator:   true,
			},
		},
	}
}

var enschede, hengelo = entity.City{Name: "Enschede"}, entity.City{Name: "Hengelo"}

func Test_User_IsValid(t *testing.T) {
	a := assert.New(t)
	testUser := getUser()
	a.Nil(testUser.IsValid())
}

func Test_User_IsValid_Invalid(t *testing.T) {
	a := assert.New(t)
	testUser := getUser()
	testUser.HousingPreferences = nil
	a.Error(testUser.IsValid())
}

func Test_User_MatchCriteria_Age(t *testing.T) {
	a := assert.New(t)
	testUser := getUser()
	testOffer := getOffer()

	testOffer.MinAge = 55
	a.False(testUser.MatchCriteria(testOffer))
	testOffer.MinAge = 18
	testOffer.MaxAge = 99
	testUser.FamilySize = 2
	testOffer.MaxIncome = 0
	a.True(testUser.MatchCriteria(testOffer))
}

func Test_User_MatchCriteria_FamilySize(t *testing.T) {
	a := assert.New(t)
	testUser := getUser()
	testOffer := getOffer()

	testOffer.MinAge = 0
	testOffer.MaxAge = 0
	a.False(testUser.MatchCriteria(testOffer))
	testUser.FamilySize = 2
	testOffer.MaxIncome = 0
	a.True(testUser.MatchCriteria(testOffer))
}

func Test_User_MatchCriteria_Income(t *testing.T) {
	a := assert.New(t)
	testUser := getUser()
	testOffer := getOffer()

	testOffer.MinAge = 0
	testOffer.MaxAge = 0
	testUser.FamilySize = 2
	a.False(testUser.MatchCriteria(testOffer))
	testOffer.MaxIncome = 40000
	a.True(testUser.MatchCriteria(testOffer))
}

func Test_User_MatchPreferences_Location(t *testing.T) {
	a := assert.New(t)
	testUser := getUser()
	testOffer := getOffer()

	a.True(testUser.MatchPreferences(testOffer))
	testUser.HousingPreferences[0].City = []entity.City{enschede, hengelo}
	a.False(testUser.MatchPreferences(testOffer))
	testOffer.Housing.City = hengelo
	a.True(testUser.MatchPreferences(testOffer))
	testOffer.Housing.City = enschede
	testUser.HousingPreferences[0].CityDistrict = []entity.CityDistrict{
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
	testOffer.Housing.CityDistrict = entity.CityDistrict{CityName: "Enschede", Name: "Enschede - Roombeek"}
	a.True(testUser.MatchPreferences(testOffer))
	testOffer.Housing.CityDistrict = entity.CityDistrict{CityName: "Enschede", Name: "deppenbroek"}
	a.False(testUser.MatchPreferences(testOffer))
}

func Test_User_MatchPreferences_HousingType(t *testing.T) {
	a := assert.New(t)
	testUser := getUser()
	testOffer := getOffer()

	a.True(testUser.MatchPreferences(testOffer))
	testOffer.Housing.Type = entity.HousingType{
		Type: entity.HousingTypeUndefined,
	}
	a.False(testUser.MatchPreferences(testOffer))
	testUser.HousingPreferences[0].Type = nil
	a.True(testUser.MatchPreferences(testOffer))
}

func Test_User_MatchPreferences_Price(t *testing.T) {
	a := assert.New(t)
	testUser := getUser()
	testOffer := getOffer()

	a.True(testUser.MatchPreferences(testOffer))
	testOffer.Housing.Price = 1000
	a.False(testUser.MatchPreferences(testOffer))
}

func Test_User_MatchPreferences_Criteria(t *testing.T) {
	a := assert.New(t)
	testUser := getUser()
	testOffer := getOffer()

	a.True(testUser.MatchPreferences(testOffer))
	testOffer.Housing.Garden = true
	testUser.HousingPreferences[0].HasGarden = false
	a.True(testUser.MatchPreferences(testOffer))
	testUser.HousingPreferences[0].NumberBedroom = 5
	a.False(testUser.MatchPreferences(testOffer))
}
