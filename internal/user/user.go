package user

import (
	"time"

	"gorm.io/gorm"

	"github.com/woningfinder/woningfinder/internal/corporation"
)

// User defines an user of WoningFinder
type User struct {
	gorm.Model
	FullName           string
	BirthDate          time.Time
	HousingPreferences HousingPreferences
}

// HousingPreferences defines the user preference on a housing
type HousingPreferences struct {
	gorm.Model
	UserID                     int                       `gorm:"primaryKey"`
	Type                       []corporation.HousingType `gorm:"many2many:housing_preferences_housing_types"`
	MinimumPrice, MaximumPrice float64
	City                       []corporation.City `gorm:"many2many:housing_preferences_cities"`
	NumberBedroom              int
	HasBalcony                 bool
	HasGarage                  bool
	HasGarden                  bool
	HasElevator                bool
	HasHousingAllowance        bool
	IsAccessible               bool
}
