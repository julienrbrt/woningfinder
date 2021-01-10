package entity

import (
	"gorm.io/gorm"
)

// HousingPreferences defines the user preference on a housing
type HousingPreferences struct {
	gorm.Model
	UserID              int           `gorm:"primaryKey"`
	Type                []HousingType `gorm:"many2many:housing_preferences_housing_types"`
	MaximumPrice        float64
	City                []City         `gorm:"many2many:housing_preferences_cities"`
	CityDistrict        []CityDistrict `gorm:"many2many:housing_preferences_city_districts"`
	NumberBedroom       int
	HasBalcony          bool
	HasGarage           bool
	HasGarden           bool
	HasElevator         bool
	HasHousingAllowance bool
	HasAttic            bool
	IsAccessible        bool
	OfferMatch          []OfferMatch
}
