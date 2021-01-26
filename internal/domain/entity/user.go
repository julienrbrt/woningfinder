package entity

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

// User defines an user of WoningFinder
type User struct {
	ID                      uint
	CreatedAt               time.Time `pg:"default:now()"`
	DeletedAt               time.Time `json:"-"`
	Name                    string
	Email                   string `pg:",unique"`
	BirthYear               int
	YearlyIncome            int
	FamilySize              int
	Plan                    UserPlan                  `pg:"rel:has-one,fk:id"`
	HousingPreferences      []HousingPreferences      `pg:"rel:has-many,join_fk:user_id"`
	HousingPreferencesMatch []HousingPreferencesMatch `pg:"rel:has-many,join_fk:user_id"`
	CorporationCredentials  []CorporationCredentials  `pg:"rel:has-many,join_fk:user_id"`
}

// Bind permits go-chi router to verify the user input and marshal it
func (u *User) Bind(r *http.Request) error {
	return u.IsValid()
}

// Render permits go-chi router to render the user
func (*User) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// IsValid checks if a user is valid
func (u *User) IsValid() error {
	if u.Name == "" || u.Email == "" {
		return fmt.Errorf("user name or email missing")
	}

	if u.BirthYear < 1900 {
		return fmt.Errorf("user birth year invalid")
	}

	if u.YearlyIncome < -1 {
		return fmt.Errorf("user yearly income invalid")
	}

	if len(u.HousingPreferences) == 0 {
		return fmt.Errorf("user must have a housing preferences")
	}

	// verify plan
	if !u.Plan.Name.Exists() {
		return fmt.Errorf("user plan invalid")
	}

	if len(u.HousingPreferences) > u.Plan.Name.MaxHousingPreferences() {
		return fmt.Errorf("error cannot create more than %d housing preferences in plan %s: got %d", u.Plan.Name.MaxHousingPreferences(), u.Plan.Name, len(u.HousingPreferences))
	}

	return nil
}

// HasPaid checks if a user has a paid plan
func (u *User) HasPaid() bool {
	if u.Plan == (UserPlan{}) || !u.Plan.Name.Exists() {
		return false
	}

	return u.Plan.CreatedAt != (time.Time{})
}

// MatchCriteria verifies that an user match the offer criterias
func (u *User) MatchCriteria(offer Offer) bool {
	age := time.Now().Year() - u.BirthYear
	// checks if offer age is set and check boundaries
	if offer.MinAge > 0 && ((age < offer.MinAge) || (offer.MaxAge != 0 && age > offer.MaxAge)) {
		return false
	}

	// checks if offer family size is set and check boundaries
	if offer.MinFamilySize > 0 && (u.FamilySize < offer.MinFamilySize || (offer.MaxFamilySize > 0 && u.FamilySize > offer.MaxFamilySize)) {
		return false
	}

	// checks if offer incomes is set and check boundaries
	if offer.MinIncome > 0 && u.YearlyIncome > -1 && (u.YearlyIncome < offer.MinIncome || (offer.MaxIncome > 0 && u.YearlyIncome > offer.MaxIncome)) {
		return false
	}

	return true
}

// MatchPreferences verifies that an offer match the user preferences
func (u *User) MatchPreferences(offer Offer) bool {
	for _, pref := range u.HousingPreferences {
		// match price
		if offer.Housing.Price >= pref.MaximumPrice {
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

		// match characteristics
		if (pref.NumberBedroom > 0 && pref.NumberBedroom > offer.Housing.NumberBedroom) ||
			(pref.HasBalcony && !offer.Housing.Balcony) ||
			(pref.HasGarden && !offer.Housing.Garden) ||
			(pref.HasElevator && !offer.Housing.Elevator) ||
			(pref.HasHousingAllowance && !offer.Housing.HousingAllowance) ||
			(pref.IsAccessible && !offer.Housing.Accessible) ||
			(pref.HasGarage && !offer.Housing.Garage) ||
			(pref.HasAttic && !offer.Housing.Attic) {
			continue
		}

		return true
	}

	return false
}

func matchHouseType(pref HousingPreferences, housing Housing) bool {
	if len(pref.Type) == 0 {
		return true
	}

	for _, t := range pref.Type {
		if t.Type == housing.Type.Type {
			return true
		}
	}

	return false
}

func matchCity(pref HousingPreferences, housing Housing) bool {
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

func matchCityDistrict(pref HousingPreferences, housing Housing) bool {
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
