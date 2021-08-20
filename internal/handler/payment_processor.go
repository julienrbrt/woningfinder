package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/render"
	"github.com/woningfinder/woningfinder/internal/customer"
	handlerErrors "github.com/woningfinder/woningfinder/internal/handler/errors"
	"github.com/woningfinder/woningfinder/pkg/util"
)

type PaymentMethod string

const (
	PaymentMethodCrypto = "crypto"
	PaymentMethodStripe = "stripe"
)

type paymentProcessorRequest struct {
	Email  string        `json:"email"`
	Method PaymentMethod `json:"method"`
}

// Bind permits go-chi router to verify the user input and marshal it
func (p *paymentProcessorRequest) Bind(r *http.Request) error {
	if !util.IsEmailValid(p.Email) {
		return errors.New("please enter a valid email")
	}

	if p.Method != PaymentMethodCrypto && p.Method != PaymentMethodStripe {
		return fmt.Errorf("invalid payment method: %s", string(p.Method))
	}

	return nil
}

// Render permits go-chi router to render the user
func (*paymentProcessorRequest) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (h *handler) PaymentProcessor(w http.ResponseWriter, r *http.Request) {
	request := &paymentProcessorRequest{}
	if err := render.Bind(r, request); err != nil {
		render.Render(w, r, handlerErrors.ErrorRenderer(err))
		return
	}

	// check if user exists
	user, err := h.userService.GetUser(&customer.User{Email: request.Email})
	if err != nil {
		render.Render(w, r, handlerErrors.ErrNotFound)
		return
	}

	if user.Plan.IsValid() {
		render.Render(w, r, handlerErrors.ErrorRenderer(errors.New("user already paid")))
		return
	}

	switch request.Method {
	case PaymentMethodStripe:
		// process payment by creating a Stripe session ID
		h.createCheckoutSession(request.Email, customer.PlanFromName(user.Plan.PlanName), w, r)
	case PaymentMethodCrypto:
		// TODO in #73
	}
}
