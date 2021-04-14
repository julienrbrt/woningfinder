package customer

import (
	"errors"
	"fmt"
	"time"

	"github.com/woningfinder/woningfinder/internal/corporation"
)

// HousingPreferences defines the user preference on a housing
type HousingPreferences struct {
	ID                  uint                      `pg:",pk" json:"-"`
	CreatedAt           time.Time                 `pg:"default:now()" json:"-"`
	UserID              uint                      `json:"-"`
	Type                []corporation.HousingType `pg:"-" json:"type"` // linked to HousingPreferencesHousingType
	MaximumPrice        float64                   `json:"maximum_price"`
	City                []corporation.City        `pg:"-" json:"city,omitempty"` // linked to HousingPreferencesCity and HousingPreferencesCityDistrict
	NumberBedroom       int                       `json:"number_bedroom"`
	HasBalcony          bool                      `json:"has_balcony"`
	HasGarage           bool                      `json:"has_garage"`
	HasGarden           bool                      `json:"has_garden"`
	HasElevator         bool                      `json:"has_elevator"`
	HasAttic            bool                      `json:"has_attic"`
	HasHousingAllowance bool                      `json:"has_housing_allowance"`
	IsAccessible        bool                      `json:"is_accessible"`
}

// HasMinimal ensure that the housing preferences contains the minimal required data
func (h *HousingPreferences) HasMinimal() error {
	if len(h.Type) == 0 {
		return errors.New("error housing preferences invalid: housing type missing")
	}

	for _, t := range h.Type {
		if t != corporation.HousingTypeAppartement && t != corporation.HousingTypeHouse {
			return fmt.Errorf("error housing preferences invalid: housing type %s does not exist", t)
		}
	}

	if h.MaximumPrice <= 0 {
		return errors.New("error housing preferences invalid: maximum price must be greater than 0")
	}

	if len(h.City) == 0 {
		return errors.New("error housing preferences invalid: cities missing")
	}

	return nil
}

// HousingPreferencesMatch defines an offer that matched with an user
// It is used to determined to which offer WoningFinder has applied
type HousingPreferencesMatch struct {
	ID              uint      `pg:",pk" json:"id"`
	CreatedAt       time.Time `pg:"default:now()" json:"created_at"`
	DeletedAt       time.Time `pg:",soft_delete" json:"-"`
	UserID          uint      `json:"-"`
	HousingAddress  string    `json:"housing_address"`
	CorporationName string    `json:"corporation_name"`
	OfferURL        string    `json:"offer_url"`
}
