package customer

import (
	"math"
	"time"
)

// Plan defines the different plans name
type Plan string

const (
	// to be updated every year with the "passend toewijze"
	MaximumIncomeSocialHouse = 44655

	// PlanBasis is the plan for social houses
	PlanBasis Plan = "basis"
	// PlanPro is the free sector houses
	PlanPro Plan = "pro"
)

// Price returns the plan price in euro
func (p Plan) Price() int {
	switch p {
	case PlanBasis:
		return 10
	case PlanPro:
		return 35
	default:
		return 0
	}
}

func (p Plan) MaximumIncome() int {
	if p == PlanBasis {
		return MaximumIncomeSocialHouse
	}

	return math.MaxInt32
}

// UserPlan stores the user plan and payment details (when paid)
type UserPlan struct {
	UserID    uint      `pg:",pk" json:"-"`
	CreatedAt time.Time `pg:"default:now()" json:"created_at"`
	Name      Plan      `json:"name"`
}
