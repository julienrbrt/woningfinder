package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/render"
	stripe "github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/checkout/session"
	stripeCustomer "github.com/stripe/stripe-go/v72/customer"
	"github.com/woningfinder/woningfinder/internal/customer"
	handlerErrors "github.com/woningfinder/woningfinder/internal/handler/errors"
	"github.com/woningfinder/woningfinder/pkg/util"
)

const (
	successURL = "https://woningfinder.nl/login?thanks=true"
	cancelURL  = "https://woningfinder.nl/start/voltooien?cancelled=true"
)

type paymentRequest struct {
	Email string `json:"email"`
}

type paymentResponse struct {
	StripeSessionID string `json:"stripe_session_id"`
}

// Bind permits go-chi router to verify the user input and marshal it
func (p *paymentRequest) Bind(r *http.Request) error {
	if !util.IsEmailValid(p.Email) {
		return errors.New("please enter a valid email")
	}

	return nil
}

// Render permits go-chi router to render the user
func (*paymentRequest) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (h *handler) PaymentProcessor(w http.ResponseWriter, r *http.Request) {
	request := &paymentRequest{}
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

	if user.Plan.IsFree() || user.Plan.IsSubscribed() {
		render.Render(w, r, handlerErrors.BadRequestErrorRenderer(errors.New("user has free plan or is already subscribed")))
		return
	}

	plan, err := customer.PlanFromName(user.Plan.Name)
	if err != nil {
		render.Render(w, r, handlerErrors.ServerErrorRenderer(err))
		return
	}

	// create customer in stripe for subscription
	customer, err := h.createStripeCustomer(user)
	if err != nil {
		errorMsg := fmt.Errorf("failed creating stripe customer")
		h.logger.Sugar().Errorf("%w: %w", errorMsg, err)
		render.Render(w, r, handlerErrors.ServerErrorRenderer(err))
		return
	}

	// creating a stripe checkout session
	session, err := h.createStripeCheckoutSession(customer, plan.StripeProductID)
	if err != nil {
		errorMsg := fmt.Errorf("error while creating stripe new checkout session")
		h.logger.Sugar().Errorf("%w: %w", errorMsg, err)
		render.Render(w, r, handlerErrors.ServerErrorRenderer(err))
		return
	}

	// return response
	json.NewEncoder(w).Encode(paymentResponse{
		StripeSessionID: session.ID,
	})
}

func (h *handler) createStripeCustomer(user *customer.User) (*stripe.Customer, error) {
	// if stripe customer already exists get it
	if user.Plan.StripeCustomerID != "" {
		if customer, err := stripeCustomer.Get(user.Plan.StripeCustomerID, &stripe.CustomerParams{}); err == nil {
			return customer, nil
		}
	}

	params := &stripe.CustomerParams{
		Name:  stripe.String(user.Name),
		Email: stripe.String(user.Email),
	}

	// enrich stripe customer with metadata
	params.AddMetadata("user_id", fmt.Sprint(user.ID))

	// create new stripe customer
	return stripeCustomer.New(params)
}

func (h *handler) createStripeCheckoutSession(customer *stripe.Customer, priceID string) (*stripe.CheckoutSession, error) {
	params := &stripe.CheckoutSessionParams{
		Customer: stripe.String(customer.ID),
		PaymentMethodTypes: stripe.StringSlice([]string{
			"card",
		}),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String(priceID),
				Quantity: stripe.Int64(1),
			},
		},
		Locale:     stripe.String("nl"),
		Mode:       stripe.String(string(stripe.CheckoutSessionModeSubscription)),
		SuccessURL: stripe.String(successURL),
		CancelURL:  stripe.String(cancelURL),
	}

	return session.New(params)
}
