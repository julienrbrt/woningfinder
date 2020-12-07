package user

import (
	"time"

	"github.com/woningfinder/woningfinder/internal/corporation"
)

// User defines an user of WoningFinder
type User struct {
	FullName           string
	BirthDate          time.Time
	HousingPreferences HousingPreferences
}

// HousingPreferences defines the user preference on a housing
type HousingPreferences struct {
	Type                       []corporation.HousingType
	MinimumPrice, MaximumPrice float64
	District                   []string
	City                       []string

	NumberBedroom       int
	HasBalcony          bool
	HasGarage           bool
	HasGarden           bool
	HasElevator         bool
	HasHousingAllowance bool
	IsAccessible        bool
}
