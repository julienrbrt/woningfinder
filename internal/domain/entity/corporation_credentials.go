package entity

import "time"

// CorporationCredentials holds the user credentials to login to an housing corporation
type CorporationCredentials struct {
	CreatedAt       time.Time `pg:"default:now()"`
	UpdatedAt       time.Time
	DeletedAt       time.Time
	UserID          int         `pg:",pk"`
	CorporationName string      `pg:",unique"`
	CorporationURL  string      `pg:",unique"`
	Corporation     Corporation `pg:"rel:has-one"`
	Login           string
	Password        string
}
