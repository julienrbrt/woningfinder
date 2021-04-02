package entity

import (
	"time"
)

// Plan defines the different plans name
type Plan string

const (
	// PlanBasis is the normal plan
	PlanBasis Plan = "basis"
	// PlanPro is the high-end plan
	PlanPro Plan = "pro"
)

// Price returns the plan price in euro
func (p Plan) Price() int {
	switch p {
	case PlanBasis:
		return 35
	case PlanPro:
		return 75
	default:
		return 0
	}
}

// UserPlan stores the user plan and payment details (when paid)
type UserPlan struct {
	UserID    uint      `pg:",pk" json:"user_id"`
	CreatedAt time.Time `pg:"default:now()" json:"created_at"`
	Name      Plan      `json:"name"`
}
