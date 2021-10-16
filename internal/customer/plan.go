package customer

import (
	"fmt"
	"math"
	"time"
)

// TODO update every year with the "passend toewijze"
const MaximumIncomeSocialHouse = 44655

// Plan defines the different plans
type Plan struct {
	StripeProductID string `json:"-"`
	Name            string `json:"name"`
	Price           int    `json:"price"`
	MaximumIncome   int    `json:"maximum_income"`
}

// PlanBasis is the free plan for social houses
var PlanBasis = Plan{
	Name:          "basis",
	Price:         0,
	MaximumIncome: MaximumIncomeSocialHouse,
}

// PlanPro is the free sector houses
var PlanPro = Plan{
	StripeProductID: "price_1JlDuPHWufZqidI1zXocKAsS",
	Name:            "pro",
	Price:           15,
	MaximumIncome:   math.MaxInt32,
}

func PlanFromName(name string) (Plan, error) {
	switch name {
	case PlanBasis.Name:
		return PlanBasis, nil
	case PlanPro.Name:
		return PlanPro, nil
	case "test-ugly-woningfinder-plan": // used for tests :(
		return Plan{Price: 1000, StripeProductID: "price_1JlDZnHWufZqidI1hmzlgann"}, nil
	}

	return Plan{}, fmt.Errorf("cannot find plan from name: %s invalid", name)
}

// UserPlan stores the user plan and payment details
type UserPlan struct {
	UserID                uint      `pg:",pk" json:"-"`
	StripeCustomerID      string    `json:"-"`
	CreatedAt             time.Time `pg:"default:now()" json:"created_at"`
	ActivatedAt           time.Time `json:"activated_at"`
	SubscriptionStartedAt time.Time `json:"subscription_started_at"`
	LastPaymentSucceeded  bool      `json:"last_payment_succeeded"`
	Name                  string    `json:"name"`
}

func (u *UserPlan) IsValid() bool {
	if plan, _ := PlanFromName(u.Name); plan.Price > 0 {
		return u.IsActivated() && u.IsSubscribed()
	}

	return u.IsActivated()
}

func (u *UserPlan) IsActivated() bool {
	return u.ActivatedAt != (time.Time{})
}

func (u *UserPlan) IsFree() bool {
	if plan, _ := PlanFromName(u.Name); plan.Price == 0 {
		return true
	}

	return false
}

func (u *UserPlan) IsSubscribed() bool {
	return u.SubscriptionStartedAt != (time.Time{}) && u.LastPaymentSucceeded
}
