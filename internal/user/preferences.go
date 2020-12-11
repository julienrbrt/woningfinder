package user

import (
	"strings"

	"github.com/woningfinder/woningfinder/internal/corporation"
)

// MatchPreferences verifies that an offer match the user preferences
func (u *User) MatchPreferences(offer corporation.Offer) bool {
	// verify if not already applied
	if offer.HasApplied {
		return false
	}

	// match price
	if offer.Housing.Price <= u.HousingPreferences.MinimumPrice || offer.Housing.Price >= u.HousingPreferences.MaximumPrice {
		return false
	}

	// match house type
	if !u.matchPrefType(offer.Housing.Type) {
		return false
	}

	// match location
	offer.Housing.SetDistrict()
	if !matchPrefLocation(offer.Housing.City, u.HousingPreferences.City) {
		return false
	}

	// match characteristics
	if (u.HousingPreferences.NumberBedroom > 0 && u.HousingPreferences.NumberBedroom > offer.Housing.NumberBedroom) ||
		(u.HousingPreferences.HasBalcony && !offer.Housing.Balcony) ||
		(u.HousingPreferences.HasGarden && !offer.Housing.Garden) ||
		(u.HousingPreferences.HasElevator && !offer.Housing.Elevator) ||
		(u.HousingPreferences.HasHousingAllowance && !offer.Housing.HousingAllowance) ||
		(u.HousingPreferences.IsAccessible && (!offer.Housing.AccessibilityScooter && !offer.Housing.AccessibilityWheelchair)) ||
		(u.HousingPreferences.HasGarage && !offer.Housing.Garage) {
		return false
	}

	return true
}

func (u *User) matchPrefType(actual corporation.HousingType) bool {
	if len(u.HousingPreferences.Type) == 0 {
		return true
	}

	for _, p := range u.HousingPreferences.Type {
		if p == actual {
			return true
		}
	}

	return false
}

func matchPrefLocation(actualCity corporation.City, pref []corporation.City) bool {
	if len(pref) == 0 {
		return true
	}

	for _, p := range pref {
		// prevent that if actual is an empty, then strings.Contains returns true
		if p.District == nil && actualCity.Name != "" && strings.Contains(strings.ToLower(actualCity.Name), strings.ToLower(p.Name)) {
			return true
		}

		for _, d := range p.District {
			if len(actualCity.District) > 0 && d.Name != "" &&
				strings.Contains(strings.ToLower(actualCity.District[0].Name), strings.ToLower(d.Name)) &&
				strings.Contains(strings.ToLower(actualCity.Name), strings.ToLower(p.Name)) {
				return true
			}
		}
	}

	return false
}
