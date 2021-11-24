package user

import (
	"fmt"

	"github.com/julienrbrt/woningfinder/internal/customer"
)

func (s *service) CreateWaitingList(w *customer.WaitingList) error {
	if _, err := s.dbClient.Conn().
		Model(w).
		OnConflict("(email, city_name) DO UPDATE").
		Insert(); err != nil {
		return fmt.Errorf("error when user %s:%s to waiting list: %w", w.Email, w.CityName, err)
	}

	return nil
}
