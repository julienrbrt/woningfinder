package user

import (
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// User defines an user of WoningFinder
type User struct {
	gorm.Model
	Name                   string
	Email                  string
	BirthYear              int
	YearlyIncome           int
	FamilySize             int
	Plan                   Plan
	HousingPreferences     []HousingPreferences
	CorporationCredentials []CorporationCredentials `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (s *service) CreateUser(u *User) (*User, error) {
	if u.Email == "" {
		return nil, fmt.Errorf("email is required for creating user")
	}

	if u.Plan != Zeker && u.Plan != Sneller {
		return nil, fmt.Errorf("subscribing to an existing plan is required")
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

func (s *service) GetUser(email string) (*User, error) {
	var u User
	s.dbClient.Conn().Where(&User{Email: email}).First(&u)

	if u.ID == 0 {
		return nil, fmt.Errorf("no user found with the email: %s", email)
	}

	return &u, nil
}

func (s *service) DeleteUser(u *User) error {
	// delete all corporations credentials
	if err := s.dbClient.Conn().Unscoped().Select(clause.Associations).
		Where(&CorporationCredentials{UserID: int(u.ID)}).
		Delete(&CorporationCredentials{}).Error; err != nil {
		return fmt.Errorf("failed to delete corporation credentials associated to this user: %w", err)
	}

	// delete housing preferences
	err := s.DeleteHousingPreferences(u)
	if err != nil {
		return fmt.Errorf("failed deleting housing preferences from user: %w", err)
	}

	return s.dbClient.Conn().Unscoped().Select(clause.Associations).Delete(u).Error
}
