package customer

import "time"

// WaitingList holds the user waiting list information
type WaitingList struct {
	CreatedAt time.Time `pg:"default:now()"`
	Email     string    `pg:",pk"`
	CityName  string    `pg:",pk"`
}
