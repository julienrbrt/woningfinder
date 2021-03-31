package entity

import (
	"database/sql/driver"
	"time"
)

// Plan defines the different plans name
type Plan string

const (
	// PlanBasis is the normal plan
	PlanBasis Plan = "BASIS"
	// PlanPro is the high-end plan
	PlanPro Plan = "PRO"
)

// Scan implements the Scanner interface from reading from the database
func (p *Plan) Scan(value interface{}) error {
	if value == nil { // happens if we try to get an user without a plan (i.e. not paid yet)
		return nil
	}

	*p = Plan(string(value.([]byte)))
	return nil
}

// Value implements the Valuer interface for the storing in the database
func (p Plan) Value() (driver.Value, error) {
	return string(p), nil
}

// MaxHousingPreferences returns the maximum autorized of housing preferences
func (p Plan) MaxHousingPreferences() int {
	switch p {
	case PlanBasis:
		return 1
	case PlanPro:
		return 10
	default:
		return 0
	}
}

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
	UserID    uint      `pg:",pk"`
	CreatedAt time.Time `pg:"default:now()"`
	DeletedAt time.Time `pg:",soft_delete" json:"-"`
	Name      Plan
}
