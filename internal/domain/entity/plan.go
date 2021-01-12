package entity

import (
	"database/sql/driver"
	"time"

	"gorm.io/gorm"
)

// Plan is the database representation of a Plan
type Plan struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Name      P              `gorm:"primaryKey"`
}

const (
	// PlanZeker is the normal plan
	PlanZeker P = "ZEKER"
	// PlanSneller is the high-end plan
	PlanSneller = "SNELLER"
)

// P defines the different plans name
type P string

// Scan implements the Scanner interface from reading from the database
func (u *P) Scan(value interface{}) error {
	*u = P(value.(string))
	return nil
}

// Value implements the Valuer interface for the storing in the database
func (u P) Value() (driver.Value, error) {
	return string(u), nil
}

// Exists check if the plan exists
func (u *P) Exists() bool {
	return *u != PlanZeker && *u != PlanSneller
}

// AllowMultipleHousingPreferences checks if the plan allows multiple housing preferences
func (u *P) AllowMultipleHousingPreferences() bool {
	return string(*u) == PlanSneller
}
