package customer_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/woningfinder/woningfinder/internal/corporation"
	"github.com/woningfinder/woningfinder/internal/customer"
)

func getUser() customer.User {
	return customer.User{
		Name:         "Test",
		Email:        "test@example.org",
		BirthYear:    1990,
		YearlyIncome: 30000,
		FamilySize:   3,
		Plan: customer.UserPlan{
			Name: customer.PlanBasis,
		},
		HousingPreferences: customer.HousingPreferences{

			Type: []corporation.HousingType{
				corporation.HousingTypeHouse,
				corporation.HousingTypeAppartement,
			},
			MaximumPrice:  950,
			NumberBedroom: 1,
			HasElevator:   true,
			City: []corporation.City{
				{Name: "Enschede"},
			},
		},
	}
}

func Test_User_HasMinimal(t *testing.T) {
	a := assert.New(t)
	testUser := getUser()
	testUser.Plan = customer.UserPlan{CreatedAt: time.Now(), Name: customer.PlanBasis}
	a.Nil(testUser.HasMinimal())
}

func Test_User_HasMinimal_InvalidPlan(t *testing.T) {
	a := assert.New(t)
	testUser := getUser()
	testUser.Plan = customer.UserPlan{CreatedAt: time.Now(), Name: "invalid"}
	a.Error(testUser.HasMinimal())
}

func Test_User_HasPaid(t *testing.T) {
	a := assert.New(t)
	testUser := getUser()
	testUser.Plan = customer.UserPlan{CreatedAt: time.Now(), Name: customer.PlanBasis}
	a.True(testUser.HasPaid())
}

func Test_User_HasPaid_Invalid(t *testing.T) {
	a := assert.New(t)
	testUser := getUser()
	a.False(testUser.HasPaid())
}
