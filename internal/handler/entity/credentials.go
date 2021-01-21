package entity

import (
	"errors"
	"net/http"
)

// CredentialsRequest permits to update and save the corporation credentials
type CredentialsRequest struct {
	CorporationName string
	Login, Password string
}

// Bind permits go-chi router to verify the user input and marshal it
func (c *CredentialsRequest) Bind(r *http.Request) error {
	if c.CorporationName == "" || c.Login == "" || c.Password == "" {
		return errors.New("given credentials are invalid")
	}

	return nil
}

// Render permits go-chi router to render the user
func (*CredentialsRequest) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// CredentialsResponse is used to display which housing corporation are supported for the user housing preferences
type CredentialsResponse struct {
	CorporationName string
	IsKnown         bool
}
