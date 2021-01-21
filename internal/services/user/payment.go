package user

import (
	"fmt"
	"time"

	"github.com/woningfinder/woningfinder/internal/domain/entity"
)

func (s *service) SetPaid(u *entity.User, plan entity.Plan) error {
	if _, err := s.dbClient.Conn().
		Model(&entity.UserPlan{UserID: u.ID, PaymentDate: time.Now(), Name: plan}).
		OnConflict("(user_id) DO UPDATE").
		Insert(); err != nil {
		return fmt.Errorf("error when adding user plan: %w", err)
	}

	return nil
}

// TODO to test
// SetExpired set a user to expired when he found a house
func (s *service) SetExpired(u *entity.User) error {
	if _, err := s.dbClient.Conn().
		Model(&entity.UserPlan{UserID: u.ID, PaymentDate: time.Time{}}).
		WherePK().
		UpdateNotZero(); err != nil {
		return fmt.Errorf("error when setting user plan expiration: %w", err)
	}

	return nil
}
