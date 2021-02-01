package matcher_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/woningfinder/woningfinder/internal/domain/entity"
	"github.com/woningfinder/woningfinder/internal/domain/matcher"
)

func getUser() *entity.User {
	return &entity.User{
		Name:         "Test",
		Email:        "test@example.org",
		BirthYear:    1990,
		YearlyIncome: 30000,
		FamilySize:   3,
		Plan: entity.UserPlan{
			Name: entity.PlanBasis,
		},
		HousingPreferences: []entity.HousingPreferences{
			{
				Type: []entity.HousingType{
					{Type: entity.HousingTypeHouse},
					{Type: entity.HousingTypeAppartement},
				},
				MaximumPrice:  950,
				NumberBedroom: 1,
				HasElevator:   true,
			},
		},
	}
}

func getOffer() entity.Offer {
	return entity.Offer{
		ExternalID: "w758752",
		Housing: entity.Housing{
			Type: entity.HousingType{
				Type: entity.HousingTypeHouse,
			},
			Latitude:      52.133,
			Longitude:     6.61433,
			Address:       "Beatrixstraat 1 R 7161 DJ Neede A",
			EnergieLabel:  "A",
			Price:         656.39,
			Size:          80,
			NumberRoom:    6,
			NumberBedroom: 2,
			BuildingYear:  2010,
			Garden:        false,
			Garage:        false,
			Elevator:      true,
			Balcony:       true,
			Attic:         false,
			Accessible:    true,
		},
	}
}

var enschede, hengelo = entity.City{Name: "Enschede"}, entity.City{Name: "Hengelo"}

func Test_MatchCriteria_Age(t *testing.T) {
	a := assert.New(t)
	testUser := getUser()
	testOffer := getOffer()

	testOffer.MinAge = 55
	a.False(matcher.MatchCriteria(testUser, testOffer))
	testOffer.MinAge = 18
	a.True(matcher.MatchCriteria(testUser, testOffer))
}

func Test_MatchCriteria_FamilySize(t *testing.T) {
	a := assert.New(t)
	testUser := getUser()
	testOffer := getOffer()

	a.True(matcher.MatchCriteria(testUser, testOffer))
	testOffer.MaxFamilySize = 2
	a.False(matcher.MatchCriteria(testUser, testOffer))
}

func Test_MatchCriteria_PassendToewijzen(t *testing.T) {
	a := assert.New(t)
	testUser := getUser()
	testOffer := getOffer()

	a.True(matcher.MatchCriteria(testUser, testOffer))
	testUser.YearlyIncome = 40000
	a.False(matcher.MatchCriteria(testUser, testOffer))
}
