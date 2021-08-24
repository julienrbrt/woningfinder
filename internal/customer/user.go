package customer

import (
	"fmt"
	"net/http"
	"time"

	"github.com/woningfinder/woningfinder/pkg/util"
)

// User defines an user of WoningFinder
type User struct {
	ID           uint      `json:"-"`
	CreatedAt    time.Time `pg:"default:now()" json:"created_at,omitempty"`
	DeletedAt    time.Time `json:"-"`
	Name         string    `json:"name"`
	Email        string    `pg:",unique" json:"email"`
	BirthYear    int       `json:"birth_year"`
	YearlyIncome int       `json:"yearly_income"`
	FamilySize   int       `json:"family_size"`
	// only used when the user is less than 23, housing allowance depends on age
	HasChildrenSameHousing  bool                      `json:"has_children_same_housing"`
	Plan                    UserPlan                  `pg:"rel:has-one,fk:id" json:"plan,omitempty"`
	HousingPreferences      HousingPreferences        `pg:"rel:has-one,fk:id" json:"housing_preferences,omitempty"`
	HousingPreferencesMatch []HousingPreferencesMatch `pg:"rel:has-many,join_fk:user_id" json:"housing_preferences_match,omitempty"`
	CorporationCredentials  []CorporationCredentials  `pg:"rel:has-many,join_fk:user_id" json:"corporation_credentials,omitempty"`
}

// HasMinimal ensure that the user contains the minimal required data
func (u *User) HasMinimal() error {
	if u.Name == "" {
		return fmt.Errorf("user name missing")
	}

	if !util.IsEmailValid(u.Email) {
		return fmt.Errorf("user email invalid")
	}

	if u.BirthYear < 1920 || u.BirthYear >= time.Now().Year() {
		return fmt.Errorf("user birth year invalid")
	}

	if u.YearlyIncome < -1 {
		return fmt.Errorf("user yearly income invalid")
	}

	plan, err := PlanFromName(u.Plan.PlanName)
	if err != nil {
		return err
	}

	if u.YearlyIncome > plan.MaximumIncome {
		return fmt.Errorf("error plan %s not allowed: yearly incomes too high", plan.Name)
	}

	if err := u.HousingPreferences.HasMinimal(); err != nil {
		return fmt.Errorf("user housing preferences invalid: %w", err)
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
