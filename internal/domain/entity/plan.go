package entity

import (
	"database/sql/driver"
	"time"
)

const (
	// PlanZeker is the normal plan
	PlanZeker Plan = "ZEKER"
	// PlanSneller is the high-end plan
	PlanSneller = "SNELLER"
)

// UserPlan stores the user plan and payment details (when paid)
// TODO when a user found a house reset paymentdate
type UserPlan struct {
	UserID      int       `pg:",pk"`
	CreatedAt   time.Time `pg:"default:now()"`
	UpdatedAt   time.Time
	DeletedAt   time.Time `pg:",soft_delete" json:"-"`
	PaymentDate time.Time `json:"-"` // When payment date is set the user has paid and the search can start
	Name        Plan
}

// Plan defines the different plans name
type Plan string

// Scan implements the Scanner interface from reading from the database
func (p *Plan) Scan(value interface{}) error {
	*p = Plan(string(value.([]byte)))
	return nil
}

// Value implements the Valuer interface for the storing in the database
func (p Plan) Value() (driver.Value, error) {
	return string(p), nil
}

// Exists check if the plan exists
func (p Plan) Exists() bool {
	return p == PlanZeker || p == PlanSneller
}

// MaxHousingPreferences returns the maximum autorized of housing preferences
func (p Plan) MaxHousingPreferences() int {
	switch p {
	case PlanZeker:
		return 1
	case PlanSneller:
		return 10
	default:
		return 0
	}
}
