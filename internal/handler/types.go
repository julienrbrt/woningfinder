package handler

import (
	"errors"
	"net/http"
)

// credentialsRequest permits to update and save the corporation credentials
type credentialsRequest struct {
	CorporationName string
	Login, Password string
}

// Bind permits go-chi router to verify the user input and marshal it
func (c *credentialsRequest) Bind(r *http.Request) error {
	if c.CorporationName == "" || c.Login == "" || c.Password == "" {
		return errors.New("given credentials are invalid")
	}

	return nil
}

// Render permits go-chi router to render the user
func (*credentialsRequest) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// credentialsResponse is used to display which housing corporation are supported for the user housing preferences
type credentialsResponse struct {
	CorporationName string
	IsKnown         bool
}
