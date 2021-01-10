package user

import (
	"fmt"

	"github.com/woningfinder/woningfinder/internal/domain/entity"
	"gorm.io/gorm/clause"
)

func (s *service) CreateUser(u *entity.User) (*entity.User, error) {
	if u.Email == "" {
		return nil, fmt.Errorf("email is required for creating user")
	}

	if !u.Plan.Exists() {
		return nil, fmt.Errorf("subscribing to a valid plan is required")
	}

	if u.YearlyIncome < -1 {
		return nil, fmt.Errorf("yearly income must be greater than 0, or set to -1 to not be used")
	}

	if len(u.HousingPreferences) == 0 {
		return nil, fmt.Errorf("housing preferences is required for creating user")
	}

	_, err := s.GetUser(u.Email)
	if err == nil {
		return nil, fmt.Errorf("error user %s already exists", u.Email)
	}

	if err := s.dbClient.Conn().Create(&u).Error; err != nil {
		return nil, err
	}

	if err := s.CreateHousingPreferences(u, u.HousingPreferences); err != nil {
		return nil, fmt.Errorf("error when creating user %s: %w", u.Email, err)
	}

	return u, nil
}

func (s *service) GetUser(email string) (*entity.User, error) {
	var u entity.User
	s.dbClient.Conn().Where(&entity.User{Email: email}).First(&u)

	if u.ID == 0 {
		return nil, fmt.Errorf("no user found with the email: %s", email)
	}

	return &u, nil
}

func (s *service) DeleteUser(u *entity.User) error {
	// delete all corporations credentials
	if err := s.dbClient.Conn().Unscoped().Select(clause.Associations).
		Where(&entity.CorporationCredentials{UserID: int(u.ID)}).
		Delete(&entity.CorporationCredentials{}).Error; err != nil {
		return fmt.Errorf("failed to delete corporation credentials associated to this user: %w", err)
	}

	// delete housing preferences
	err := s.DeleteHousingPreferences(u)
	if err != nil {
		return fmt.Errorf("failed deleting housing preferences from user: %w", err)
	}

	return s.dbClient.Conn().Unscoped().Select(clause.Associations).Delete(u).Error
}
