package matcher_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/woningfinder/woningfinder/internal/domain/entity"
	"github.com/woningfinder/woningfinder/internal/domain/matcher"
)

func Test_MatchPreferences_Location(t *testing.T) {
	a := assert.New(t)
	testUser := getUser()
	testOffer := getOffer()

	a.True(matcher.MatchPreferences(testUser, testOffer))
	testUser.HousingPreferences[0].City = []entity.City{enschede, hengelo}
	a.False(matcher.MatchPreferences(testUser, testOffer))
	testOffer.Housing.City = hengelo
	a.True(matcher.MatchPreferences(testUser, testOffer))
	testOffer.Housing.City = enschede
	testUser.HousingPreferences[0].CityDistrict = []entity.HousingPreferencesCityDistrict{
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
	a.False(matcher.MatchPreferences(testUser, testOffer))
	testOffer.Housing.CityDistrict = "Enschede - Roombeek"
	a.True(matcher.MatchPreferences(testUser, testOffer))
	testOffer.Housing.CityDistrict = "deppenbroek"
	a.False(matcher.MatchPreferences(testUser, testOffer))
}

func Test_MatchPreferences_HousingType(t *testing.T) {
	a := assert.New(t)
	testUser := getUser()
	testOffer := getOffer()

	a.True(matcher.MatchPreferences(testUser, testOffer))
	testOffer.Housing.Type = entity.HousingType{
		Type: entity.HousingTypeUndefined,
	}
	a.False(matcher.MatchPreferences(testUser, testOffer))
	testUser.HousingPreferences[0].Type = nil
	a.True(matcher.MatchPreferences(testUser, testOffer))
}

func Test_MatchPreferences_Price(t *testing.T) {
	a := assert.New(t)
	testUser := getUser()
	testOffer := getOffer()

	a.True(matcher.MatchPreferences(testUser, testOffer))
	testOffer.Housing.Price = 1000
	a.False(matcher.MatchPreferences(testUser, testOffer))
}

func Test_MatchPreferences_HousingAllowance(t *testing.T) {
	a := assert.New(t)
	testUser := getUser()
	testOffer := getOffer()

	a.True(matcher.MatchPreferences(testUser, testOffer))
	testOffer.Housing.Price = 1000
	testUser.HousingPreferences[0].MaximumPrice = 1000
	testUser.HousingPreferences[0].HasHousingAllowance = false
	a.True(matcher.MatchPreferences(testUser, testOffer))
	testUser.HousingPreferences[0].HasHousingAllowance = true
	a.False(matcher.MatchPreferences(testUser, testOffer))
}

func Test_MatchPreferences_Criteria(t *testing.T) {
	a := assert.New(t)
	testUser := getUser()
	testOffer := getOffer()

	a.True(matcher.MatchPreferences(testUser, testOffer))
	testOffer.Housing.Garden = true
	testUser.HousingPreferences[0].HasGarden = false
	a.True(matcher.MatchPreferences(testUser, testOffer))
	testUser.HousingPreferences[0].NumberBedroom = 5
	a.False(matcher.MatchPreferences(testUser, testOffer))
}
