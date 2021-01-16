package entity

import (
	"fmt"
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
	City                []City                           `pg:"many2many:housing_preferences_cities" json:",omitempty"`
	CityDistrict        []HousingPreferencesCityDistrict `pg:"rel:has-many" json:",omitempty"`
	NumberBedroom       int
	HasBalcony          bool
	HasGarage           bool
	HasGarden           bool
	HasElevator         bool
	HasHousingAllowance bool
	HasAttic            bool
	IsAccessible        bool
}

// HousingPreferencesCityDistrict defines the user preferences city districts
type HousingPreferencesCityDistrict struct {
	HousingPreferencesID int
	Name                 string
	CityName             string
}

// IsValid verifies that the given HousingPreferences is valid
func (h *HousingPreferences) IsValid() error {
	if len(h.Type) == 0 {
		return fmt.Errorf("housing preferences housing type missing")
	}

	if len(h.City) == 0 {
		return fmt.Errorf("housing preferences cities missing")
	}

	return nil
}

// HousingPreferencesMatch defines an offer that matched with an user
// It is used to determined to which offer WoningFinder has applied
type HousingPreferencesMatch struct {
	ID              int       `pg:",pk"`
	CreatedAt       time.Time `pg:"default:now()"`
	UpdatedAt       time.Time
	DeletedAt       time.Time `pg:",soft_delete"`
	UserID          int
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
