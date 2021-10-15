package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/render"
	stripe "github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/checkout/session"
	"github.com/woningfinder/woningfinder/internal/customer"
	handlerErrors "github.com/woningfinder/woningfinder/internal/handler/errors"
	"github.com/woningfinder/woningfinder/pkg/util"
)

const (
	successURL = "https://woningfinder.nl/login?thanks=true"
	cancelURL  = "https://woningfinder.nl/start/voltooien?cancelled=true"
)

type subscriptionRequest struct {
	Email string `json:"email"`
}

type subscriptionResponse struct {
	StripeSessionID string `json:"stripe_session_id"`
}

// Bind permits go-chi router to verify the user input and marshal it
func (p *subscriptionRequest) Bind(r *http.Request) error {
	if !util.IsEmailValid(p.Email) {
		return errors.New("please enter a valid email")
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

	plan, err := customer.PlanFromName(user.Plan.Name)
	if err != nil {
		render.Render(w, r, handlerErrors.ServerErrorRenderer(err))
		return
	}

	// process payment by creating a Stripe session ID
	h.createStripeCheckoutSession(request.Email, plan, w, r)
}

func (h *handler) createStripeCheckoutSession(email string, plan customer.Plan, w http.ResponseWriter, r *http.Request) {
	params := &stripe.CheckoutSessionParams{
		CustomerEmail: stripe.String(email),
		SubmitType:    stripe.String("pay"),
		PaymentMethodTypes: stripe.StringSlice([]string{
			"ideal", "card",
		}),
		PaymentIntentData: &stripe.CheckoutSessionPaymentIntentDataParams{
			ReceiptEmail: stripe.String(email),
		},
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			planToLineItems(plan),
		},
		Locale:     stripe.String("nl"),
		Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL: stripe.String(successURL),
		CancelURL:  stripe.String(cancelURL),
	}

	session, err := session.New(params)
	if err != nil {
		errorMsg := fmt.Errorf("error while creating stripe new checkout session")
		h.logger.Sugar().Errorf("%w: %w", errorMsg, err)
		render.Render(w, r, handlerErrors.ServerErrorRenderer(errorMsg))
		return
	}

	// return response
	json.NewEncoder(w).Encode(subscriptionResponse{
		StripeSessionID: session.ID,
	})
}

// planToLineItems gets the plan costs and converts it to cents
func planToLineItems(plan customer.Plan) *stripe.CheckoutSessionLineItemParams {
	switch plan {
	case customer.PlanPro:
		return &stripe.CheckoutSessionLineItemParams{
			Currency: stripe.String(string(stripe.CurrencyEUR)),
			Name:     stripe.String(strings.Title(customer.PlanPro.Name)),
			Amount:   stripe.Int64(int64(customer.PlanPro.Price) * 100),
			Quantity: stripe.Int64(1),
		}
	}

	return nil
}
