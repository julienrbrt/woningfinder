package matcher

import (
	"strings"
	"time"

	"github.com/woningfinder/woningfinder/internal/corporation"
	"github.com/woningfinder/woningfinder/internal/corporation/city"
	"github.com/woningfinder/woningfinder/internal/customer"
	"github.com/woningfinder/woningfinder/pkg/logging"
)

type Matcher interface {
	MatchOffer(user customer.User, offer corporation.Offer) bool
}

type matcher struct{}

func NewMatcher() Matcher {
	return &matcher{}
}

func (m *matcher) MatchOffer(user customer.User, offer corporation.Offer) bool {
	return m.matchCriteria(user, offer) && m.matchPreferences(user.HousingPreferences, offer)
}

// matchCriteria verifies that an user match the offer criterias
func (m *matcher) matchCriteria(user customer.User, offer corporation.Offer) bool {
	age := time.Now().Year() - user.BirthYear

	// checks if offer age is set and check boundaries
	if offer.MinAge > 0 && ((age < offer.MinAge) || (offer.MaxAge > 0 && age > offer.MaxAge)) {
		return false
	}

	// checks if offer family size is set and check boundaries
	if offer.MinFamilySize > 0 && (user.FamilySize < offer.MinFamilySize) || (offer.MaxFamilySize > 0 && user.FamilySize > offer.MaxFamilySize) {
		return false
	}

	// checks if offer incomes is set and check boundaries
	min, max := m.passendToewijzen(user)
	if offer.Housing.Price < min || offer.Housing.Price > max {
		return false
	}

	return true
}

// matchPreferences verifies that an offer match the user preferences
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

	for _, preferences := range housingPreference.City {
		// prevent that if actual is an empty, then strings.Contains returns true
		if housing.City.Name != "" && strings.EqualFold(housing.City.Name, preferences.Name) {
			if len(preferences.District) > 0 {
				return matchDistrict(preferences, housing)
			}

			return true
		}
	}

	return false
}

func matchDistrict(cityPreferences city.City, housing corporation.Housing) bool {
	// no district but user has preferences so reject
	if housing.CityDistrict == "" {
		return false
	}

	cityDistricts := city.SuggestedCityDistrictFromName(&logging.Logger{}, cityPreferences.Name)
	for district := range cityPreferences.District {
		// get neighboorhoud for district and try finding a match
		if neighboorhoud, ok := cityDistricts[district]; ok {
			for _, n := range neighboorhoud {
				if strings.EqualFold(n, housing.CityDistrict) {
					return true
				}
			}
		}

		// advanced user defined district immediatly so perform basic check
		if strings.Contains(strings.ToLower(district), strings.ToLower(housing.CityDistrict)) || strings.Contains(strings.ToLower(housing.CityDistrict), strings.ToLower(district)) {
			return true
		}
	}

	return false
}
