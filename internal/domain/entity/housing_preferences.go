package entity

import (
	"time"
)

// HousingPreferences defines the user preference on a housing
type HousingPreferences struct {
	ID                  int       `pg:",pk"`
	CreatedAt           time.Time `pg:"default:now()"`
	UpdatedAt           time.Time
	DeletedAt           time.Time `pg:",soft_delete"`
	UserID              int
	Type                []HousingType `pg:"many2many:housing_preferences_housing_types,join_fk:housing_type"`
	MaximumPrice        float64
	City                []City         `pg:"many2many:housing_preferences_cities"`
	CityDistrict        []CityDistrict `pg:"many2many:housing_preferences_city_districts"`
	NumberBedroom       int
	HasBalcony          bool
	HasGarage           bool
	HasGarden           bool
	HasElevator         bool
	HasHousingAllowance bool
	HasAttic            bool
	IsAccessible        bool
}

// HousingPreferencesMatch defines an offer that matched with an user
// It is used to determined to which offer WoningFinder has applied
type HousingPreferencesMatch struct {
	ID              int       `pg:",pk"`
	CreatedAt       time.Time `pg:"default:now()"`
	UpdatedAt       time.Time
	DeletedAt       time.Time `pg:",soft_delete"`
	UserID          int       `pg:",pk"`
	HousingAddress  string
	CorporationName string
	OfferURL        string
}

// HousingPreferencesHousingType defines the many-to-many relationship table
type HousingPreferencesHousingType struct {
	HousingPreferencesID int
	HousingType          Type
}

// HousingPreferencesCity defines the many-to-many relationship table
type HousingPreferencesCity struct {
	HousingPreferencesID int
	CityName             string
}

// HousingPreferencesCityDistrict defines the many-to-many relationship table
type HousingPreferencesCityDistrict struct {
	HousingPreferencesID int
	CityName             string
	CityDistrictName     string
}
