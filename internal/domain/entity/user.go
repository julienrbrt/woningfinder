package entity

import (
	"fmt"
	"net/http"
	"time"
)

// User defines an user of WoningFinder
type User struct {
	ID           uint      `json:"id"`
	CreatedAt    time.Time `pg:"default:now()" json:"created_at"`
	DeletedAt    time.Time `json:"-"`
	Name         string    `json:"name"`
	Email        string    `pg:",unique" json:"email"`
	BirthYear    int       `json:"birth_year"`
	YearlyIncome int       `json:"yearly_income"`
	FamilySize   int       `json:"family_size"`
	// only used when the user is less than 23, housing allowance depends on age
	HasChildrenSameHousing  bool                      `json:"has_children_same_housing"`
	Plan                    UserPlan                  `pg:"rel:has-one,fk:id" json:"plan,omitempty"`
	HousingPreferences      []HousingPreferences      `pg:"rel:has-many,join_fk:user_id" json:"housing_preferences,omitempty"`
	HousingPreferencesMatch []HousingPreferencesMatch `pg:"rel:has-many,join_fk:user_id" json:"housing_prefereces_match,omitempty"`
	CorporationCredentials  []CorporationCredentials  `pg:"rel:has-many,join_fk:user_id" json:"corporation_credentials,omitempty"`
}

// HasMinimal ensure that the user contains the minimal required data
func (u *User) HasMinimal() error {
	if u.Name == "" || u.Email == "" {
		return fmt.Errorf("user name or email missing")
	}

	if u.BirthYear < 1900 || u.BirthYear >= time.Now().Year() {
		return fmt.Errorf("user birth year invalid")
	}

	if u.YearlyIncome < -1 {
		return fmt.Errorf("user yearly income invalid")
	}

	if len(u.HousingPreferences) == 0 {
		return fmt.Errorf("user must have a housing preferences")
	}

	// verify housing preferences
	if len(u.HousingPreferences) > u.Plan.Name.MaxHousingPreferences() {
		return fmt.Errorf("error cannot create more than %d housing preferences in plan %s: got %d", u.Plan.Name.MaxHousingPreferences(), u.Plan.Name, len(u.HousingPreferences))
	}

	for _, housingPreferences := range u.HousingPreferences {
		if err := housingPreferences.HasMinimal(); err != nil {
			return fmt.Errorf("user housing preferences invalid: %w", err)
		}
	}

	return nil
}

// Bind permits go-chi router to verify the user input and marshal it
func (u *User) Bind(r *http.Request) error {
	return u.HasMinimal()
}

// Render permits go-chi router to render the user
func (*User) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// HasPaid checks if a user has a paid plan
func (u *User) HasPaid() bool {
	if u.Plan == (UserPlan{}) {
		return false
	}

	return u.Plan.CreatedAt != (time.Time{})
}
