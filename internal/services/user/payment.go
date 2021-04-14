package user

import (
	"fmt"

	"github.com/woningfinder/woningfinder/internal/customer"
)

func (s *service) SetPaid(userID uint, plan customer.Plan) error {
	if _, err := s.dbClient.Conn().
		Model(&customer.UserPlan{UserID: userID, Name: plan}).
		OnConflict("(user_id) DO UPDATE").
		Insert(); err != nil {
		return fmt.Errorf("error when adding user plan: %w", err)
	}

	return nil
}
