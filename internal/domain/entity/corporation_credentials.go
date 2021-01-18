package entity

import (
	"fmt"
	"time"
)

// CorporationCredentials holds the user credentials to login to an housing corporation
type CorporationCredentials struct {
	CreatedAt       time.Time `pg:"default:now()"`
	UpdatedAt       time.Time
	DeletedAt       time.Time   `json:"-"`
	UserID          int         `pg:",pk"`
	CorporationName string      `pg:",pk"`
	Corporation     Corporation `pg:"rel:has-one"`
	Login           string
	Password        string
}

// IsValid verifies the validity of the corporation credentials
func (c *CorporationCredentials) IsValid() error {
	if c.Corporation.Name == "" {
		return fmt.Errorf("corporation invalid")
	}

	if c.Login == "" || c.Password == "" {
		return fmt.Errorf("login or password missing")
	}

	return nil
}
