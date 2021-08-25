package customer_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/woningfinder/woningfinder/internal/corporation"
	"github.com/woningfinder/woningfinder/internal/corporation/city"
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
			PlanName: customer.PlanBasis.Name,
		},
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
}

func Test_User_HasMinimal(t *testing.T) {
	a := assert.New(t)
	testUser := getUser()
	testUser.Plan = customer.UserPlan{CreatedAt: time.Now(), PlanName: customer.PlanBasis.Name}
	a.Nil(testUser.HasMinimal())
}

func Test_User_HasMinimal_InvalidPlan(t *testing.T) {
	a := assert.New(t)
	testUser := getUser()
	testUser.Plan = customer.UserPlan{CreatedAt: time.Now(), PlanName: "invalid"}
	a.Error(testUser.HasMinimal())
}

func Test_User_Plan(t *testing.T) {
	a := assert.New(t)
	testUser := getUser()
	a.False(testUser.Plan.IsValid())
	a.False(testUser.Plan.IsFreeTrialValid())
	a.False(testUser.Plan.IsPaid())
	testUser.Plan = customer.UserPlan{}
	a.False(testUser.Plan.IsValid())
	a.False(testUser.Plan.IsFreeTrialValid())
	a.False(testUser.Plan.IsPaid())
	testUser.Plan = customer.UserPlan{CreatedAt: time.Date(2020, 01, 01, 0, 0, 0, 0, time.UTC), PlanName: customer.PlanBasis.Name}
	a.False(testUser.Plan.IsValid())
	a.False(testUser.Plan.IsPaid())
	a.False(testUser.Plan.IsFreeTrialValid())
	testUser.Plan = customer.UserPlan{CreatedAt: time.Now(), PlanName: customer.PlanBasis.Name}
	a.True(testUser.Plan.IsValid())
	a.False(testUser.Plan.IsPaid())
	a.True(testUser.Plan.IsFreeTrialValid())
	testUser.Plan = customer.UserPlan{CreatedAt: time.Date(2020, 01, 01, 0, 0, 0, 0, time.UTC), PlanName: customer.PlanBasis.Name, PurchasedAt: time.Now()}
	a.True(testUser.Plan.IsValid())
	a.True(testUser.Plan.IsPaid())
	a.False(testUser.Plan.IsFreeTrialValid())
}
