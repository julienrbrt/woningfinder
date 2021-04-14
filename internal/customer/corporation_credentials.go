package customer

import (
	"errors"
	"net/http"
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

// Bind permits go-chi router to verify the user input and marshal it
func (c *CorporationCredentials) Bind(r *http.Request) error {
	if c.CorporationName == "" || c.Login == "" || c.Password == "" {
		return errors.New("credentials cannot be empty")
	}

	return nil
}

// Render permits go-chi router to render the user
func (*CorporationCredentials) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
