package user

import (
	"strings"

	"github.com/woningfinder/woningfinder/internal/corporation"
)

// MatchPreferences verifies that an offer match the user preferences
func (u *User) MatchPreferences(offer corporation.Offer) bool {
	// match price
	if offer.Housing.Price >= u.HousingPreferences.MaximumPrice {
		return false
	}

	// match house type
	if !u.matchHouseType(offer.Housing) {
		return false
	}

	// match location
	if !u.matchCity(offer.Housing) || !u.matchCityDistrict(offer.Housing) {
		return false
	}

	// match characteristics
	if (u.HousingPreferences.NumberBedroom > 0 && u.HousingPreferences.NumberBedroom > offer.Housing.NumberBedroom) ||
		(u.HousingPreferences.HasBalcony && !offer.Housing.Balcony) ||
		(u.HousingPreferences.HasGarden && !offer.Housing.Garden) ||
		(u.HousingPreferences.HasElevator && !offer.Housing.Elevator) ||
		(u.HousingPreferences.HasHousingAllowance && !offer.Housing.HousingAllowance) ||
		(u.HousingPreferences.IsAccessible && !offer.Housing.Accessible) ||
		(u.HousingPreferences.HasGarage && !offer.Housing.Garage) ||
		(u.HousingPreferences.HasAttic && !offer.Housing.Attic) {
		return false
	}

	return true
}

func (u *User) matchHouseType(housing corporation.Housing) bool {
	if len(u.HousingPreferences.Type) == 0 {
		return true
	}

	for _, t := range u.HousingPreferences.Type {
		if t.Type == housing.Type.Type {
			return true
		}
	}

	return false
}

func (u *User) matchCity(housing corporation.Housing) bool {
	if len(u.HousingPreferences.City) == 0 {
		return true
	}

	for _, city := range u.HousingPreferences.City {
		// prevent that if actual is an empty, then strings.Contains returns true
		if housing.City.Name != "" && strings.Contains(strings.ToLower(housing.City.Name), strings.ToLower(city.Name)) {
			return true
		}
	}

	return false
}

func (u *User) matchCityDistrict(housing corporation.Housing) bool {
	// no preferences so accept everything
	if len(u.HousingPreferences.CityDistrict) == 0 {
		return true
	}

	// no district but user has preferences so reject
	if housing.CityDistrict.Name == "" {
		return false
	}

	for _, district := range u.HousingPreferences.CityDistrict {
		// prevent that if actual is an empty, then strings.Contains returns true
		if housing.CityDistrict.CityName != "" && strings.Contains(strings.ToLower(housing.CityDistrict.CityName), strings.ToLower(district.CityName)) &&
			housing.CityDistrict.Name != "" && strings.Contains(strings.ToLower(housing.CityDistrict.Name), strings.ToLower(district.Name)) {
			return true
		}
	}

	return false
}
