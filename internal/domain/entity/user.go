package entity

import (
	"fmt"
	"net/http"
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
	HasChildrenSameHousing  bool                      // only used when the user is less than 23, the housing allowance depends then
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

	if len(u.HousingPreferences) > u.Plan.Name.MaxHousingPreferences() {
		return fmt.Errorf("error cannot create more than %d housing preferences in plan %s: got %d", u.Plan.Name.MaxHousingPreferences(), u.Plan.Name, len(u.HousingPreferences))
	}

	return nil
}

// HasPaid checks if a user has a paid plan
func (u *User) HasPaid() bool {
	if u.Plan == (UserPlan{}) {
		return false
	}

	return u.Plan.CreatedAt != (time.Time{})
}
