package user

import (
	"errors"
	"fmt"
	"strings"
	"time"

	pg "github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/julienrbrt/woningfinder/internal/customer"
	"go.uber.org/zap"
)

var ErrUserAlreadyExist = errors.New("user already exist")

// CreateUser creates an user
func (s *service) CreateUser(user *customer.User) error {
	db := s.dbClient.Conn()

	// verify user
	if err := user.HasMinimal(); err != nil {
		return fmt.Errorf("error user %s invalid: %w", user.Email, err)
	}

	// a user cannot be activated at creation
	user.ActivatedAt = (time.Time{})

	// activate failed reaction email alert by default
	user.HasAlertsEnabled = true

	// create user - if exist throw error
	if _, err := db.Model(user).Insert(); err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint \"users_email_key\"") {
			return ErrUserAlreadyExist
		}

		return fmt.Errorf("failing creating user: %w", err)
	}

	// create user housing preferences
	if err := s.CreateHousingPreferences(user.ID, &user.HousingPreferences); err != nil {
		// rollback
		if err2 := s.DeleteUser(user.Email); err2 != nil {
			s.logger.Error("error when rolling back user creation: %w", zap.Any("originalErr", err), zap.Error(err2))
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

// ConfirmUser validate user account
func (s *service) ConfirmUser(email string) error {
	if _, err := s.dbClient.Conn().
		Model((*customer.User)(nil)).
		Set("activated_at = now()").
		Where("email ILIKE ?", email).
		Update(); err != nil {
		return fmt.Errorf("error when activating user: %w", err)
	}

	return nil
}

// UpdateUser update the user basics info (name, revenues and family size)
func (s *service) UpdateUser(user *customer.User) error {
	if _, err := s.dbClient.Conn().
		Model((*customer.User)(nil)).
		Set("name = ?", user.Name).
		Set("family_size = ?", user.FamilySize).
		Set("yearly_income = ?", user.YearlyIncome).
		Set("has_alerts_enabled = ?", user.HasAlertsEnabled).
		Where("id = ?", user.ID).
		Update(); err != nil {
		return fmt.Errorf("error when updating user informaion: %w", err)
	}

	return nil
}

// GetUsersWithGivenCorporationCredentials gets all the users with a given corporation credentials
func (s *service) GetUsersWithGivenCorporationCredentials(corporationName string) ([]*customer.User, error) {
	var users []*customer.User
	if err := s.dbClient.Conn().
		Model(&users).
		Relation("CorporationCredentials", func(q *orm.Query) (*orm.Query, error) {
			return q.Where("corporation_name = ?", corporationName), nil
		}).
		Join("INNER JOIN corporation_credentials cc ON id = cc.user_id").
		Where("cc.corporation_name = ?", corporationName).
		Order("created_at ASC"). // first created user first
		Select(); err != nil {
		return nil, fmt.Errorf("error getting users with %s corporation credentials: %w", corporationName, err)
	}

	return users, nil
}

// GetWeeklyUpdateUsers gets all reactions of users from the last week
func (s *service) GetWeeklyUpdateUsers() ([]*customer.User, error) {
	var users []*customer.User
	if err := s.dbClient.Conn().
		Model(&users).
		Relation("CorporationCredentials").
		Relation("HousingPreferencesMatch", func(q *orm.Query) (*orm.Query, error) {
			return q.Where("created_at >= now() - interval '7 day'"), nil
		}).
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
	if _, err := db.Model((*customer.User)(nil)).Where("id = ?", user.ID).Delete(); err != nil && !errors.Is(err, pg.ErrNoRows) {
		return fmt.Errorf("failed delete user for %s: %w", email, err)
	}

	if _, err := db.Model((*customer.HousingPreferencesMatch)(nil)).Where("user_id = ?", user.ID).Delete(); err != nil && !errors.Is(err, pg.ErrNoRows) {
		return fmt.Errorf("failed deleting housing preferences match for %s: %w", email, err)
	}

	return nil
}
