package user

import (
	"fmt"

	"github.com/woningfinder/woningfinder/internal/domain/entity"
)

// TODO eventually use a prepare function to create it in one query only
func (s *service) CreateUser(u *entity.User) error {
	db := s.dbClient.Conn()

	// verify user
	if err := u.IsValid(); err != nil {
		return fmt.Errorf("error user %s invalid: %w", u.Email, err)
	}

	// associate to tier
	var tier entity.Tier
	if err := db.Model(&tier).Where("name = ?", u.Tier.Name).Select(); err != nil {
		return fmt.Errorf("error getting tier for creating user: %w", err)
	}
	u.Tier.ID = tier.ID

	// create user - if exist throw error
	if _, err := db.Model(&u).Insert(); err != nil {
		return fmt.Errorf("failing creating user: %w", err)
	}

	// create user housing preferences
	if err := s.CreateHousingPreferences(u, u.HousingPreferences); err != nil {
		return fmt.Errorf("error when creating user %s: %w", u.Email, err)
	}

	return nil
}

func (s *service) GetUser(email string) (*entity.User, error) {
	var user entity.User
	if err := db.Model(&user).Where("email = ?", email).Select(); err != nil {
		return nil, fmt.Errorf("failed getting user %s: %w", email, err)
	}

	if u.ID == 0 {
		return nil, fmt.Errorf("no user found with the email: %s", email)
	}

	// note enriching the user is not necessary because GetUser is almost never used

	return &u, nil
}

func (s *service) DeleteUser(u *entity.User) error {
	// TODO to implement
	// delete all corporations credentials
	// delete housing preferences
	// delete user
	panic("not implemented")
}
