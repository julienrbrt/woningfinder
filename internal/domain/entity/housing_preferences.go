package entity

import (
	"errors"
	"time"
)

// HousingPreferences defines the user preference on a housing
type HousingPreferences struct {
	ID                  uint                             `pg:",pk" json:"id"`
	CreatedAt           time.Time                        `pg:"default:now()" json:"created_at"`
	DeletedAt           time.Time                        `pg:",soft_delete" json:"-"`
	UserID              uint                             `json:"user_id"`
	Type                []HousingType                    `pg:"many2many:housing_preferences_housing_types,join_fk:housing_type" json:"type"`
	MaximumPrice        float64                          `json:"maximum_price"`
	City                []City                           `pg:"many2many:housing_preferences_cities" json:",omitempty"`
	CityDistrict        []HousingPreferencesCityDistrict `pg:"rel:has-many" json:",omitempty"`
	NumberBedroom       int                              `json:"number_bedroom"`
	HasBalcony          bool                             `json:"has_balcony"`
	HasGarage           bool                             `json:"has_garage"`
	HasGarden           bool                             `json:"has_garden"`
	HasElevator         bool                             `json:"has_elevator"`
	HasAttic            bool                             `json:"has_attic"`
	HasHousingAllowance bool                             `json:"has_housing_allowance"`
	IsAccessible        bool                             `json:"is_accessible"`
}

// HasMinimal ensure that the housing preferences contains the minimal required data
func (h *HousingPreferences) HasMinimal() error {
	if len(h.Type) == 0 {
		return errors.New("error housing preferences invalid: housing type missing")
	}

	if len(h.City) == 0 {
		return errors.New("error housing preferences invalid: cities missing")
	}

	return nil
}

// HousingPreferencesCityDistrict defines the user preferences city districts
type HousingPreferencesCityDistrict struct {
	HousingPreferencesID uint
	Name                 string
	CityName             string
}

// HousingPreferencesMatch defines an offer that matched with an user
// It is used to determined to which offer WoningFinder has applied
type HousingPreferencesMatch struct {
	ID              uint      `pg:",pk" json:"id"`
	CreatedAt       time.Time `pg:"default:now()" json:"created_at"`
	DeletedAt       time.Time `pg:",soft_delete" json:"-"`
	UserID          uint      `json:"user_id"`
	HousingAddress  string    `json:"housing_address"`
	CorporationName string    `json:"corporation_name"`
	OfferURL        string    `json:"offer_url"`
}

// HousingPreferencesHousingType defines the many-to-many relationship table
type HousingPreferencesHousingType struct {
	HousingPreferencesID uint
	HousingType          string
}

// HousingPreferencesCity defines the many-to-many relationship table
type HousingPreferencesCity struct {
	HousingPreferencesID uint
	CityName             string
}
