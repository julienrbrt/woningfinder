package matcher

import (
	"strings"

	"github.com/woningfinder/woningfinder/internal/corporation"
	"github.com/woningfinder/woningfinder/internal/customer"
)

// MatchPreferences verifies that an offer match the user preferences
func (m *matcher) matchPreferences(preferences customer.HousingPreferences, offer corporation.Offer) bool {
	// match price
	if preferences.MaximumPrice > 0 && offer.Housing.Price > preferences.MaximumPrice {
		return false
	}

	// match house type
	if !matchHouseType(preferences, offer.Housing) {
		return false
	}

	// match location
	if !matchLocation(preferences, offer.Housing) {
		return false
	}

	// match characteristics
	if (preferences.NumberBedroom > 0 && preferences.NumberBedroom > offer.Housing.NumberBedroom) ||
		(preferences.IsAccessible && !offer.Housing.Accessible) ||
		(preferences.HasGarage && !offer.Housing.Garage) {
		return false
	}

	// appartement specific
	if offer.Housing.Type == corporation.HousingTypeAppartement &&
		(preferences.HasBalcony && !offer.Housing.Balcony) ||
		(preferences.HasElevator && !offer.Housing.Elevator) {
		return false
	}

	// house specific
	if offer.Housing.Type == corporation.HousingTypeHouse &&
		(preferences.HasGarden && !offer.Housing.Garden) ||
		(preferences.HasAttic && !offer.Housing.Attic) {
		return false
	}

	return true
}

func matchHouseType(housingPreference customer.HousingPreferences, housing corporation.Housing) bool {
	if len(housingPreference.Type) == 0 {
		return true
	}

	for _, t := range housingPreference.Type {
		if t == housing.Type {
			return true
		}
	}

	return false
}

func matchLocation(housingPreference customer.HousingPreferences, housing corporation.Housing) bool {
	if len(housingPreference.City) == 0 {
		return true
	}

	for _, cityPreferences := range housingPreference.City {
		// prevent that if actual is an empty, then strings.Contains returns true
		if housing.City.Name != "" && strings.Contains(strings.ToLower(housing.City.Name), strings.ToLower(cityPreferences.Name)) {
			if len(cityPreferences.District) > 0 {
				// no district but user has preferences so reject
				if housing.CityDistrict == "" {
					return false
				}

				for _, district := range cityPreferences.District {
					if strings.Contains(strings.ToLower(district), strings.ToLower(housing.CityDistrict)) || strings.Contains(strings.ToLower(housing.CityDistrict), strings.ToLower(district)) {
						return true
					}
				}

				return false
			}

			return true
		}
	}

	return false
}
