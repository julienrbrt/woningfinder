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
		if !matchCity(pref, offer.Housing) || !matchCityDistrict(pref, offer.Housing) {
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

func matchHouseType(pref entity.HousingPreferences, housing entity.Housing) bool {
	if len(pref.Type) == 0 {
		return true
	}

	for _, t := range pref.Type {
		if t == housing.Type {
			return true
		}
	}

	return false
}

func matchCity(pref entity.HousingPreferences, housing entity.Housing) bool {
	if len(pref.City) == 0 {
		return true
	}

	for _, city := range pref.City {
		// prevent that if actual is an empty, then strings.Contains returns true
		if housing.City.Name != "" && strings.Contains(strings.ToLower(housing.City.Name), strings.ToLower(city.Name)) {
			return true
		}
	}

	return false
}

func matchCityDistrict(pref entity.HousingPreferences, housing entity.Housing) bool {
	// no preferences so accept everything
	if len(pref.CityDistrict) == 0 {
		return true
	}

	// no district but user has preferences so reject
	if housing.CityDistrict == "" {
		return false
	}

	for _, district := range pref.CityDistrict {
		// prevent that if actual is an empty, then strings.Contains returns true
		if housing.City.Name != "" && strings.Contains(strings.ToLower(housing.City.Name), strings.ToLower(district.CityName)) &&
			strings.Contains(strings.ToLower(housing.CityDistrict), strings.ToLower(district.Name)) {
			return true
		}
	}

	return false
}
