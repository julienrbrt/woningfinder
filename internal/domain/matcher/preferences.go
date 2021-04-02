package matcher

import (
	"strings"

	"github.com/woningfinder/woningfinder/internal/domain/entity"
)

// MatchPreferences verifies that an offer match the user preferences
func MatchPreferences(u *entity.User, offer entity.Offer) bool {
	// match price
	if offer.Housing.Price > u.HousingPreferences.MaximumPrice {
		return false
	}

	// match house type
	if !matchHouseType(u.HousingPreferences, offer.Housing) {
		return false
	}

	// match location
	if !matchLocation(u.HousingPreferences, offer.Housing) {
		return false
	}

	// verify housing allowance
	if u.HousingPreferences.HasHousingAllowance && offer.Housing.Price > maximalehuurgrens {
		return false
	}

	// match characteristics
	if (u.HousingPreferences.NumberBedroom > 0 && u.HousingPreferences.NumberBedroom > offer.Housing.NumberBedroom) ||
		(u.HousingPreferences.HasBalcony && !offer.Housing.Balcony) ||
		(u.HousingPreferences.HasGarden && !offer.Housing.Garden) ||
		(u.HousingPreferences.HasElevator && !offer.Housing.Elevator) ||
		(u.HousingPreferences.IsAccessible && !offer.Housing.Accessible) ||
		(u.HousingPreferences.HasGarage && !offer.Housing.Garage) ||
		(u.HousingPreferences.HasAttic && !offer.Housing.Attic) {
		return false
	}

	return true
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
					if strings.Contains(strings.ToLower(housing.CityDistrict), strings.ToLower(district.Name)) {
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
