package customer

import (
	"math"
	"time"
)

const (
	// TODO update every year with the "passend toewijze"
	MaximumIncomeSocialHouse = 44655

	// Free trial of 14 days
	freeTrialDuration = 14 * 24 * time.Hour
)

// Plan defines the different plans
type Plan struct {
	Name          string `json:"name"`
	Price         int    `json:"price"`
	MaximumIncome int    `json:"maximum_income"`
}

// PlanBasis is the plan for social houses
var PlanBasis = Plan{
	Name:          "basis",
	Price:         10,
	MaximumIncome: MaximumIncomeSocialHouse,
}

// PlanPro is the free sector houses
var PlanPro = Plan{
	Name:          "pro",
	Price:         25,
	MaximumIncome: math.MaxInt32,
}

func PlanFromName(name string) Plan {
	switch name {
	case "basis":
		return PlanBasis
	case "pro":
		return PlanPro
	default:
		return Plan{}
	}
}

func PlanFromPrice(price int64) Plan {
	switch price {
	case int64(PlanBasis.Price):
		return PlanBasis
	case int64(PlanPro.Price):
		return PlanPro
	default:
		return Plan{}
	}
}

// UserPlan stores the user plan and payment details (when paid)
type UserPlan struct {
	UserID      uint      `pg:",pk" json:"-"`
	CreatedAt   time.Time `pg:"default:now()" json:"created_at"`
	PurchasedAt time.Time `json:"purchased_at"`
	PlanName    string    `json:"name"`
}

// IsValid checks if a user has a paid plan or is within its free trial
func (u *UserPlan) IsValid() bool {
	return u.PurchasedAt != (time.Time{}) || time.Until(u.CreatedAt.Add(freeTrialDuration)) > 0
}
