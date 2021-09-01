package user

import (
	"errors"
	"fmt"
	"strings"
	"time"

	pg "github.com/go-pg/pg/v10"
	"github.com/woningfinder/woningfinder/internal/customer"
)

var ErrUserAlreadyExist = errors.New("user already exist")

// CreateUser creates an user
func (s *service) CreateUser(user *customer.User) error {
	db := s.dbClient.Conn()

	// verify user
	if err := user.HasMinimal(); err != nil {
		return fmt.Errorf("error user %s invalid: %w", user.Email, err)
	}

	// a user cannot have paid when being created so reset by security
	user.Plan.FreeTrialStartedAt = (time.Time{})
	user.Plan.PurchasedAt = (time.Time{})

	// create user - if exist throw error
	if _, err := db.Model(user).Insert(); err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint \"users_email_key\"") {
			return ErrUserAlreadyExist
		}

		return fmt.Errorf("failing creating user: %w", err)
	}

	// create user plan
	if _, err := s.dbClient.Conn().Model(&customer.UserPlan{UserID: user.ID, Name: user.Plan.Name}).Insert(); err != nil {
		// rollback
		if _, err2 := db.Model(user).Where("email ILIKE ?", user.Email).Delete(); err2 != nil {
			s.logger.Sugar().Errorf("error %w and error when rolling back user creation: %w", err, err2)
		}

		return fmt.Errorf("error when creating user plan: %w", err)
	}

	// create user housing preferences
	if err := s.CreateHousingPreferences(user.ID, &user.HousingPreferences); err != nil {
		// rollback
		if _, err2 := db.Model(user).Where("email ILIKE ?", user.Email).Delete(); err2 != nil {
			s.logger.Sugar().Errorf("error %w and error when rolling back user creation: %w", err, err2)
		}

		return fmt.Errorf("error when creating user %s: %w", user.Email, err)
	}

	return nil
}

func (s *service) GetUser(email string) (*customer.User, error) {
	db := s.dbClient.Conn()

	var user customer.User
	if err := db.Model(&user).
		Where("email ILIKE ?", email).
		Relation("Plan").
		Relation("HousingPreferencesMatch").
		Select(); err != nil {
		return nil, fmt.Errorf("failed getting user %s: %w", email, err)
	}

	// enrich housing preferences
	var err error
	user.HousingPreferences, err = s.GetHousingPreferences(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed getting user %s: %w", email, err)
	}

	// enriching with corporation credentials is skipped because not useful
	return &user, nil
}

// ConfirmUser validate user account and start free trial
func (s *service) ConfirmUser(email string) error {
	// get user
	var user customer.User
	if err := s.dbClient.Conn().Model(&user).Where("email ILIKE ?", email).Select(); err != nil {
		return fmt.Errorf("failed getting user %s: %w", email, err)
	}

	// set that has started its free trial
	if _, err := s.dbClient.Conn().
		Model((*customer.UserPlan)(nil)).
		Set("free_trial_started_at = now()").
		Where("user_id = ?", user.ID).
		Update(); err != nil {
		return fmt.Errorf("error when updating user plan: %w", err)
	}

	return nil
}

// ConfirmPayment set that the user has paid
func (s *service) ConfirmPayment(email string) (*customer.User, error) {
	// get user
	var user *customer.User
	if err := s.dbClient.Conn().Model(&user).Where("email ILIKE ?", email).Select(); err != nil {
		return nil, fmt.Errorf("failed getting user %s: %w", email, err)
	}

	// set that user has paid
	if _, err := s.dbClient.Conn().
		Model((*customer.UserPlan)(nil)).
		Set("purchased_at = now()").
		Where("user_id = ?", user.ID).
		Update(); err != nil {
		return nil, fmt.Errorf("error when updating user plan: %w", err)
	}

	return user, nil
}

// GetUsersWithGivenCorporationCredentials gets all the users with a given corporation credentials
func (s *service) GetUsersWithGivenCorporationCredentials(corporationName string) ([]customer.User, error) {
	var users []customer.User
	if err := s.dbClient.Conn().
		Model(&users).
		Relation("Plan").
		Relation("CorporationCredentials").
		Join("INNER JOIN corporation_credentials cc ON id = cc.user_id").
		Where("cc.corporation_name = ?", corporationName).
		Order("created_at ASC"). // first created user first
		Select(); err != nil {
		return nil, fmt.Errorf("error getting users with %s corporation credentials: %w", corporationName, err)
	}

	return users, nil
}

// GetUsersWithHousingPreferencesMatch gets all reactions of paid user for the last week
func (s *service) GetUsersWithHousingPreferencesMatch() ([]customer.User, error) {
	var users []customer.User
	if err := s.dbClient.Conn().
		Model(&users).
		Relation("Plan").
		Relation("CorporationCredentials").
		Relation("HousingPreferencesMatch").
		Join("INNER JOIN housing_preferences_matches hpm ON \"user\".\"id\" = hpm.user_id").
		Where("hpm.created_at >= now() - interval '7 day'").
		Order("created_at ASC"). // first created user first
		Select(); err != nil {
		return nil, fmt.Errorf("failed getting users with housing preferences match: %w", err)
	}

	return users, nil
}

// DeleteUser deletes an user
func (s *service) DeleteUser(email string) error {
	db := s.dbClient.Conn()

	// get user
	var user customer.User
	if err := db.Model(&user).Where("email ILIKE ?", email).Select(); err != nil && !errors.Is(err, pg.ErrNoRows) {
		return fmt.Errorf("failed getting user %s: %w", email, err)
	}

	// delete all corporations credentials
	if _, err := db.Model((*customer.CorporationCredentials)(nil)).Where("user_id = ?", user.ID).Delete(); err != nil && !errors.Is(err, pg.ErrNoRows) {
		return fmt.Errorf("failed deleting housing corporation credentials for %s: %w", email, err)
	}

	// delete housing preferences
	if err := s.DeleteHousingPreferences(user.ID); err != nil {
		return fmt.Errorf("failed deleting housing preferences for %s: %w", email, err)
	}

	// delete user
	if _, err := db.Model((*customer.UserPlan)(nil)).Where("user_id = ?", user.ID).Delete(); err != nil && !errors.Is(err, pg.ErrNoRows) {
		return fmt.Errorf("failed delete user plan for %s: %w", email, err)
	}

	if _, err := db.Model((*customer.User)(nil)).Where("id = ?", user.ID).Delete(); err != nil && !errors.Is(err, pg.ErrNoRows) {
		return fmt.Errorf("failed delete user for %s: %w", email, err)
	}

	if _, err := db.Model((*customer.HousingPreferencesMatch)(nil)).Where("user_id = ?", user.ID).Delete(); err != nil && !errors.Is(err, pg.ErrNoRows) {
		return fmt.Errorf("failed deleting housing preferences match for %s: %w", email, err)
	}

	return nil
}
