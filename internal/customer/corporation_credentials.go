package customer

import (
	"time"
)

// CorporationCredentials holds the user credentials to login to an housing corporation
type CorporationCredentials struct {
	CreatedAt       time.Time `pg:"default:now()" json:"created_at,omitempty"`
	UserID          uint      `pg:",pk" json:"-"`
	CorporationName string    `pg:",pk" json:"corporation_name"`
	Login           string    `json:"login"`
	Password        string    `json:"password"`
	FailureCount    int       `json:"-"` // FailureCount measures the number of login failure
}
