package user

import (
	"fmt"

	"github.com/julienrbrt/woningfinder/internal/customer"
)

func (s *service) GetReminderCount(email string) (*customer.ReminderCounter, error) {
	reminder := &customer.ReminderCounter{Email: email}

	_, err := s.dbClient.Conn().
		Model(reminder).
		Where("email = ?", email).
		SelectOrInsert()
	if err != nil {
		return nil, fmt.Errorf("error when getting user reminders: %w", err)
	}

	// wether the reminder was created or not, we return it
	// if it was created, the default values will be used
	return reminder, nil
}

func (s *service) UpdateReminderCount(email string, reminder *customer.ReminderCounter) error {
	if _, err := s.dbClient.Conn().
		Model((*customer.ReminderCounter)(nil)).
		Set("corporation_credentials_missing_reminder_count = ?", reminder.CorporationCredentialsMissingReminderCount).
		Set("unconfirmed_reminder_count = ?", reminder.UnconfirmedReminderCount).
		Where("email ILIKE ?", email).
		Update(); err != nil {
		return fmt.Errorf("error when updating user reminders: %w", err)
	}

	return nil
}
