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

	// change city preferences
	testUser.HousingPreferences.City = []entity.City{{Name: "Neede"}}
	a.False(matcher.MatchPreferences(testUser, testOffer))

	// add district preferences
	testUser.HousingPreferences.City = []entity.City{
		{
			Name: "Enschede",
			District: []entity.CityDistrict{
				{
					Name: "roombeek",
				},
				{
					Name: "city (oude markt)",
				},
				{
					Name: "ribbelt - stokhorst",
				},
			},
		},
	}
	a.False(matcher.MatchPreferences(testUser, testOffer))

	// add housing city district
	testOffer.Housing.CityDistrict = "city"
	a.True(matcher.MatchPreferences(testUser, testOffer))
}

func Test_MatchPreferences_HousingType(t *testing.T) {
	a := assert.New(t)
	testUser := getUser()
	testOffer := getOffer()

	a.True(matcher.MatchPreferences(testUser, testOffer))
	testOffer.Housing.Type = entity.HousingTypeUndefined
	a.False(matcher.MatchPreferences(testUser, testOffer))
	testUser.HousingPreferences.Type = nil
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
	testUser.HousingPreferences.MaximumPrice = 1000
	testUser.HousingPreferences.HasHousingAllowance = false
	a.True(matcher.MatchPreferences(testUser, testOffer))
	testUser.HousingPreferences.HasHousingAllowance = true
	a.False(matcher.MatchPreferences(testUser, testOffer))
}

func Test_MatchPreferences_Criteria(t *testing.T) {
	a := assert.New(t)
	testUser := getUser()
	testOffer := getOffer()

	a.True(matcher.MatchPreferences(testUser, testOffer))
	testOffer.Housing.Garden = true
	testUser.HousingPreferences.HasGarden = false
	a.True(matcher.MatchPreferences(testUser, testOffer))
	testUser.HousingPreferences.NumberBedroom = 5
	a.False(matcher.MatchPreferences(testUser, testOffer))
}
