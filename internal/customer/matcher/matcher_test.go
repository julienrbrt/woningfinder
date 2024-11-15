package matcher_test

import (
	"testing"

	bootstrapCorporation "github.com/julienrbrt/woningfinder/internal/bootstrap/corporation"
	"github.com/julienrbrt/woningfinder/internal/corporation"
	"github.com/julienrbrt/woningfinder/internal/corporation/city"
	"github.com/julienrbrt/woningfinder/internal/customer"
	"github.com/julienrbrt/woningfinder/internal/customer/matcher"
	"github.com/julienrbrt/woningfinder/pkg/logging"
	"github.com/stretchr/testify/assert"
)

var user = customer.User{
	Name:         "Test",
	Email:        "test@example.org",
	BirthYear:    1990,
	YearlyIncome: 30000,
	FamilySize:   3,
	HousingPreferences: customer.HousingPreferences{

		Type: []corporation.HousingType{
			corporation.HousingTypeHouse,
			corporation.HousingTypeAppartement,
		},
		MaximumPrice:  950,
		NumberBedroom: 1,
		HasElevator:   true,
		City: []city.City{
			{Name: "Enschede"},
		},
	},
}

var offer = corporation.Offer{
	ExternalID: "w758752",
	Housing: corporation.Housing{
		Type:          corporation.HousingTypeHouse,
		CityName:      city.Enschede.Name,
		CityDistrict:  "deppenbroek",
		Address:       "Beatrixstraat 1 R 7142BM Enschede",
		Price:         656.39,
		Size:          80,
		NumberBedroom: 2,
		Garden:        false,
		Garage:        false,
		Elevator:      true,
		Balcony:       true,
		Accessible:    true,
	},
}

var (
	cities = bootstrapCorporation.CreateConnectorProvider(logging.NewZapLoggerWithoutSentry(), nil).GetCities()
	m      = matcher.NewMatcher(city.NewSuggester(cities))
)

func Test_MatchOffer(t *testing.T) {
	a := assert.New(t)
	offer := offer

	matcher := m
	a.True(matcher.MatchOffer(user, offer))
}

func Test_MatchCriteria_Age(t *testing.T) {
	a := assert.New(t)
	offer := offer

	matcher := m
	a.True(matcher.MatchOffer(user, offer))
	offer.MinAge = 55
	a.False(matcher.MatchOffer(user, offer))
	offer.MinAge = 18
	a.True(matcher.MatchOffer(user, offer))
}

func Test_MatchCriteria_FamilySize(t *testing.T) {
	a := assert.New(t)
	offer := offer

	matcher := m
	a.True(matcher.MatchOffer(user, offer))
	offer.MaxFamilySize = 2
	a.False(matcher.MatchOffer(user, offer))
}

func Test_MatchCriteria_PassendToewijzen(t *testing.T) {
	a := assert.New(t)
	user := user

	matcher := m
	a.True(matcher.MatchOffer(user, offer))
	user.YearlyIncome = 40000
	a.False(matcher.MatchOffer(user, offer))
}

func Test_MatchCriteria_Incomes(t *testing.T) {
	a := assert.New(t)
	user := user
	offer := offer

	matcher := m
	offer.Housing.Price = 950
	offer.MinimumIncome = 45000
	a.False(matcher.MatchOffer(user, offer))
	user.YearlyIncome = 55000
	a.True(matcher.MatchOffer(user, offer))
	offer.MaximumIncome = 48000
	a.False(matcher.MatchOffer(user, offer))
}

func Test_MatchPreferences_Location(t *testing.T) {
	a := assert.New(t)
	user := user

	matcher := m
	a.True(matcher.MatchOffer(user, offer))

	// change city preferences
	user.HousingPreferences.City = []city.City{{Name: "Hengelo"}}
	a.False(matcher.MatchOffer(user, offer))

	// add district preferences
	user.HousingPreferences.City = []city.City{
		{
			Name: "Enschede",
			District: []string{
				"Roombeek",
				"Glanerbrug", // from suggested districts
				"Stokhorst",
			},
		},
	}
	a.False(matcher.MatchOffer(user, offer))

	// add housing city district
	offer := offer
	offer.Housing.CityDistrict = "roombeek-roomveldje"
	a.True(matcher.MatchOffer(user, offer))

	offer.Housing.CityDistrict = "glanerveld" // this wijk is inside glanerbrug
	a.True(matcher.MatchOffer(user, offer))

	offer.Housing.CityDistrict = "stroinkslanden"
	a.False(matcher.MatchOffer(user, offer))

	offer.Housing.CityDistrict = ""
	a.False(matcher.MatchOffer(user, offer))
}

func Test_MatchPreferences_HousingType(t *testing.T) {
	a := assert.New(t)

	matcher := m
	a.True(matcher.MatchOffer(user, offer))
	offer := offer
	offer.Housing.Type = corporation.HousingTypeUndefined
	a.False(matcher.MatchOffer(user, offer))
	user := user
	user.HousingPreferences.Type = nil
	a.True(matcher.MatchOffer(user, offer))
}

func Test_MatchPreferences_Price(t *testing.T) {
	a := assert.New(t)

	matcher := m
	a.True(matcher.MatchOffer(user, offer))
	offer := offer
	offer.Housing.Price = 1000
	a.False(matcher.MatchOffer(user, offer))
}

func Test_MatchPreferences_Criteria(t *testing.T) {
	a := assert.New(t)
	offer := offer
	user := user

	matcher := m
	a.True(matcher.MatchOffer(user, offer))
	offer.Housing.Garden = true
	user.HousingPreferences.HasGarden = false
	a.True(matcher.MatchOffer(user, offer))
	user.HousingPreferences.NumberBedroom = 5
	a.False(matcher.MatchOffer(user, offer))
}
