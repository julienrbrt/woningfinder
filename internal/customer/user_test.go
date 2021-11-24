package customer_test

import (
	"testing"
	"time"

	"github.com/julienrbrt/woningfinder/internal/corporation"
	"github.com/julienrbrt/woningfinder/internal/corporation/city"
	"github.com/julienrbrt/woningfinder/internal/customer"
	"github.com/stretchr/testify/assert"
)

func getUser() customer.User {
	return customer.User{
		Name:         "Test",
		Email:        "test@example.org",
		BirthYear:    1990,
		YearlyIncome: 30000,
		FamilySize:   3,
		Plan: customer.UserPlan{
			CreatedAt: time.Date(2020, 01, 01, 0, 0, 0, 0, time.UTC),
			Name:      customer.PlanBasis.Name,
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
	testUser.Plan = customer.UserPlan{CreatedAt: time.Now(), Name: customer.PlanBasis.Name}
	a.Nil(testUser.HasMinimal())
}

func Test_User_HasMinimal_InvalidPlan(t *testing.T) {
	a := assert.New(t)
	testUser := getUser()
	testUser.Plan = customer.UserPlan{CreatedAt: time.Now(), Name: "invalid"}
	a.Error(testUser.HasMinimal())
}

func Test_User_Plan(t *testing.T) {
	a := assert.New(t)
	testUser := getUser()
	a.False(testUser.Plan.IsValid())
	a.False(testUser.Plan.IsActivated())
	a.False(testUser.Plan.IsSubscribed())
	a.True(testUser.Plan.IsFree())
	testUser.Plan = customer.UserPlan{CreatedAt: time.Date(2020, 01, 01, 0, 0, 0, 0, time.UTC), Name: customer.PlanBasis.Name, ActivatedAt: time.Date(2020, 01, 01, 0, 0, 0, 0, time.UTC)}
	a.True(testUser.Plan.IsValid())
	a.False(testUser.Plan.IsSubscribed())
	a.True(testUser.Plan.IsActivated())
	a.True(testUser.Plan.IsFree())
	testUser.Plan = customer.UserPlan{CreatedAt: time.Date(2020, 01, 01, 0, 0, 0, 0, time.UTC), Name: customer.PlanPro.Name, ActivatedAt: time.Now()}
	a.False(testUser.Plan.IsValid())
	a.False(testUser.Plan.IsSubscribed())
	a.True(testUser.Plan.IsActivated())
	a.False(testUser.Plan.IsFree())
	testUser.Plan = customer.UserPlan{CreatedAt: time.Date(2020, 01, 01, 0, 0, 0, 0, time.UTC), Name: customer.PlanPro.Name, ActivatedAt: time.Date(2020, 01, 01, 0, 0, 0, 0, time.UTC), SubscriptionStartedAt: time.Now(), LastPaymentSucceeded: true}
	a.True(testUser.Plan.IsValid())
	a.True(testUser.Plan.IsSubscribed())
	a.True(testUser.Plan.IsActivated())
	a.False(testUser.Plan.IsFree())
	testUser.Plan = customer.UserPlan{CreatedAt: time.Date(2020, 01, 01, 0, 0, 0, 0, time.UTC), Name: customer.PlanPro.Name, ActivatedAt: time.Date(2020, 01, 01, 0, 0, 0, 0, time.UTC), SubscriptionStartedAt: time.Now(), LastPaymentSucceeded: false}
	a.False(testUser.Plan.IsValid())
	a.False(testUser.Plan.IsSubscribed())
	a.True(testUser.Plan.IsActivated())
	a.False(testUser.Plan.IsFree())
}
