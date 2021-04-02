package entity

import (
	"time"
)

// CorporationCredentials holds the user credentials to login to an housing corporation
type CorporationCredentials struct {
	CreatedAt       time.Time   `pg:"default:now()"`
	UserID          uint        `pg:",pk"`
	CorporationName string      `pg:",pk"`
	Corporation     Corporation `pg:"rel:has-one"`
	Login           string
	Password        string
}
