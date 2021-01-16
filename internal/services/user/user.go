package user

import (
	"fmt"

	"github.com/woningfinder/woningfinder/internal/domain/entity"
)

// TODO eventually use a prepare function to create it in one query only
func (s *service) CreateUser(u entity.User) error {
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
	u.TierID = tier.ID

	// create user - if exist throw error
	if _, err := db.Model(&u).Insert(); err != nil {
		return fmt.Errorf("failing creating user: %w", err)
	}

	// create user housing preferences
	if err := s.CreateHousingPreferences(&u, u.HousingPreferences); err != nil {
		// rollback
		if _, err2 := db.Model(&u).Where("email = ?", u.Email).Delete(); err2 != nil {
			s.logger.Sugar().Errorf("error %w and error when rolling back user creation: %w", err, err2)
		}

		return fmt.Errorf("error when creating user %s: %w", u.Email, err)
	}

	return nil
}

func (s *service) GetUser(email string) (*entity.User, error) {
	db := s.dbClient.Conn()

	// get user
	var u entity.User
	if err := db.Model(&u).Where("email ILIKE ?", email).Select(); err != nil {
		return nil, fmt.Errorf("failed getting user %s: %w", email, err)
	}

	// enrich user
	var err error
	u.HousingPreferences, err = s.GetHousingPreferences(&u)
	if err != nil {
		return nil, fmt.Errorf("failed getting user %s housing preferences: %w", u.Email, err)
	}

	if err := db.Model(&u).Where("email = ?", u.Email).Relation("Tier").Select(); err != nil {
		return nil, fmt.Errorf("failed getting user %s tier: %w", u.Email, err)
	}

	if err := db.Model(&u).Where("email = ?", u.Email).Relation("HousingPreferencesMatch").Select(); err != nil {
		return nil, fmt.Errorf("failed getting user %s housing preferences match: %w", u.Email, err)
	}

	// enriching with corporation credentials is skipped because not useful

	return &u, nil
}

func (s *service) DeleteUser(u *entity.User) error {
	// TODO to implement
	// delete all corporations credentials
	// delete housing preferences
	// delete user
	panic("not implemented")
}
