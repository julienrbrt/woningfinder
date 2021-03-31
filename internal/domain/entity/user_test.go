package entity_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/woningfinder/woningfinder/internal/domain/entity"
)

func getUser() entity.User {
	return entity.User{
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
					entity.HousingTypeHouse,
					entity.HousingTypeAppartement,
				},
				MaximumPrice:  950,
				NumberBedroom: 1,
				HasElevator:   true,
			},
		},
	}
}

func Test_User_HasPaid(t *testing.T) {
	a := assert.New(t)
	testUser := getUser()
	testUser.Plan = entity.UserPlan{CreatedAt: time.Now(), Name: entity.PlanBasis}
	a.True(testUser.HasPaid())
}

func Test_User_HasPaid_Invalid(t *testing.T) {
	a := assert.New(t)
	testUser := getUser()
	a.False(testUser.HasPaid())
}
