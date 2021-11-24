package customer

import (
	"fmt"
	"time"

	"github.com/julienrbrt/woningfinder/pkg/util"
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
	HasChildrenSameHousing  bool                       `json:"has_children_same_housing"`
	HasAlertsEnabled        bool                       `json:"has_alerts_enabled"`
	Plan                    UserPlan                   `pg:"rel:has-one,fk:id,join_fk:user_id" json:"plan,omitempty"`
	HousingPreferences      HousingPreferences         `pg:"rel:has-one,fk:id,join_fk:user_id" json:"housing_preferences,omitempty"`
	HousingPreferencesMatch []*HousingPreferencesMatch `pg:"rel:has-many,join_fk:user_id" json:"housing_preferences_match,omitempty"`
	CorporationCredentials  []*CorporationCredentials  `pg:"rel:has-many,join_fk:user_id" json:"corporation_credentials,omitempty"`
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

	if u.YearlyIncome < 0 {
		return fmt.Errorf("user yearly income invalid")
	}

	if u.FamilySize < 0 {
		return fmt.Errorf("user family size invalid")
	}

	plan, err := PlanFromName(u.Plan.Name)
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
