package customer

import (
	"time"
)

// Plan defines the different plans name
type Plan string

const (
	// PlanBasis is the plan for social houses
	PlanBasis Plan = "basis"
	// PlanPro is the free sector houses
	PlanPro Plan = "pro"
)

// Price returns the plan price in euro
func (p Plan) Price() int {
	switch p {
	case PlanBasis:
		return 2
	case PlanPro:
		return 35
	default:
		return 0
	}
}

// UserPlan stores the user plan and payment details (when paid)
type UserPlan struct {
	UserID    uint      `pg:",pk" json:"-"`
	CreatedAt time.Time `pg:"default:now()" json:"created_at"`
	Name      Plan      `json:"name"`
}
