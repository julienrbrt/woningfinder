package entity

import (
	"errors"
	"net/http"
)

// PaymentData concerns the payment data gotten from Stripe
type PaymentData struct {
	UserEmail string
	Plan      Plan
}

// Bind permits go-chi router to verify the user input and marshal it
func (p *PaymentData) Bind(r *http.Request) error {
	if p.UserEmail == "" {
		return errors.New("given payment information are invalid")
	}

	if !p.Plan.Exists() {
		return errors.New("given paid plan is invalid")
	}

	return nil
}

// Render permits go-chi router to render the user
func (*PaymentData) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
