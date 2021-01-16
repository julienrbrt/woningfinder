package entity

import (
	"database/sql/driver"
)

// Tier is the database representation of a plan
type Tier struct {
	ID    int  `pg:",pk"`
	Name  Plan `pg:",unique"`
	Price int
}

const (
	// PlanZeker is the normal plan
	PlanZeker Plan = "ZEKER"
	// PlanSneller is the high-end plan
	PlanSneller = "SNELLER"
)

// Plan defines the different plans name
type Plan string

// Scan implements the Scanner interface from reading from the database
func (u *Plan) Scan(value interface{}) error {
	*u = Plan(string(value.([]byte)))
	return nil
}

// Value implements the Valuer interface for the storing in the database
func (u Plan) Value() (driver.Value, error) {
	return string(u), nil
}

// Exists check if the plan exists
func (u *Plan) Exists() bool {
	return *u == PlanZeker || *u == PlanSneller
}

// AllowMultipleHousingPreferences checks if the plan allows multiple housing preferences
func (u *Plan) AllowMultipleHousingPreferences() bool {
	return string(*u) == PlanSneller
}
