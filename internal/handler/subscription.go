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

	successURL = "https://woningfinder.nl/login?thanks=true"
	cancelURL  = "https://woningfinder.nl/start/voltooien?cancelled=true"
)

type subscriptionRequest struct {
	Email  string        `json:"email"`
	Method PaymentMethod `json:"method"`
}

type subscriptionResponse struct {
	StripeSessionID  string `json:"stripe_session_id"`
	CryptoPaymentURL string `json:"crypto_payment_url"`
}

// Bind permits go-chi router to verify the user input and marshal it
func (p *subscriptionRequest) Bind(r *http.Request) error {
	if !util.IsEmailValid(p.Email) {
		return errors.New("please enter a valid email")
	}

	if p.Method != PaymentMethodCrypto && p.Method != PaymentMethodStripe {
		return fmt.Errorf("invalid payment method: %s", string(p.Method))
	}

	return nil
}

// Render permits go-chi router to render the user
func (*subscriptionRequest) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (h *handler) Subscription(w http.ResponseWriter, r *http.Request) {
	request := &subscriptionRequest{}
	if err := render.Bind(r, request); err != nil {
		render.Render(w, r, handlerErrors.BadRequestErrorRenderer(err))
		return
	}

	// get given user
	user, err := h.userService.GetUser(request.Email)
	if err != nil {
		render.Render(w, r, handlerErrors.ErrNotFound)
		return
	}

	if user.Plan.IsFree() {
		render.Render(w, r, handlerErrors.BadRequestErrorRenderer(errors.New("user has free plan")))
		return
	}

	if user.Plan.IsSubscribed() {
		render.Render(w, r, handlerErrors.BadRequestErrorRenderer(errors.New("user is already subscribed")))
		return
	}

	plan, err := customer.PlanFromName(user.Plan.Name)
	if err != nil {
		render.Render(w, r, handlerErrors.ServerErrorRenderer(err))
		return
	}

	switch request.Method {
	case PaymentMethodStripe:
		// process payment by creating a Stripe session ID
		h.createStripeCheckoutSession(request.Email, plan, w, r)
	case PaymentMethodCrypto:
		// process payment by creating a Crypto.com payment ID
		h.createCryptoCheckoutSession(request.Email, plan, w, r)
	}
}
