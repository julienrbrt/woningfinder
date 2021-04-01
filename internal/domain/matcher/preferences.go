package matcher

import (
	"strings"

	"github.com/woningfinder/woningfinder/internal/domain/entity"
)

// MatchPreferences verifies that an offer match the user preferences
func MatchPreferences(u *entity.User, offer entity.Offer) bool {
	for _, pref := range u.HousingPreferences {
		// match price
		if offer.Housing.Price > pref.MaximumPrice {
			continue
		}

		// match house type
		if !matchHouseType(pref, offer.Housing) {
			continue
		}

		// match location
		if !matchLocation(pref, offer.Housing) {
			continue
		}

		// verify housing allowance
		if pref.HasHousingAllowance && offer.Housing.Price > maximalehuurgrens {
			continue
		}

		// match characteristics
		if (pref.NumberBedroom > 0 && pref.NumberBedroom > offer.Housing.NumberBedroom) ||
			(pref.HasBalcony && !offer.Housing.Balcony) ||
			(pref.HasGarden && !offer.Housing.Garden) ||
			(pref.HasElevator && !offer.Housing.Elevator) ||
			(pref.IsAccessible && !offer.Housing.Accessible) ||
			(pref.HasGarage && !offer.Housing.Garage) ||
			(pref.HasAttic && !offer.Housing.Attic) {
			continue
		}

		return true
	}

	return false
}

func matchHouseType(housingPreference entity.HousingPreferences, housing entity.Housing) bool {
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

func matchLocation(housingPreference entity.HousingPreferences, housing entity.Housing) bool {
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
					if strings.Contains(strings.ToLower(housing.CityDistrict), strings.ToLower(district)) {
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
