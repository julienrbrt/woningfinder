package customer

import "time"

// WaitingList holds the user waiting list information
type WaitingList struct {
	CreatedAt time.Time `pg:"default:now()"`
	Email     string    `pg:",pk"`
	CityName  string    `pg:",pk"`
}

// ReminderCounter holds if a user already got a reminder about something
type ReminderCounter struct {
	Email                                      string `pg:",pk"`
	CorporationCredentialsMissingReminderCount int    `pg:",use_zero"`
	UnconfirmedReminderCount                   int    `pg:",use_zero"`
}
