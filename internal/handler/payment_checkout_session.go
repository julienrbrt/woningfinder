package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/render"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/checkout/session"
	"github.com/woningfinder/woningfinder/internal/customer"
	handlerErrors "github.com/woningfinder/woningfinder/internal/handler/errors"
	"github.com/woningfinder/woningfinder/pkg/util"
)

type paymentProcessorRequest struct {
	Email string        `json:"email"`
	Plan  customer.Plan `json:"plan"`
}

// Bind permits go-chi router to verify the user input and marshal it
func (p *paymentProcessorRequest) Bind(r *http.Request) error {
	if !util.IsEmailValid(p.Email) {
		return errors.New("please give a valid email")
	}

	if p.Plan.Price() == 0 {
		return fmt.Errorf("error plan %s does not exist", p.Plan)
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

	if user.HasPaid() {
		render.Render(w, r, handlerErrors.ErrorRenderer(errors.New("user already paid")))
		return
	}

	// process payment by creating a session id from Stripe
	h.createCheckoutSession(request.Email, request.Plan, w, r)
}

type createCheckoutSessionResponse struct {
	SessionID string `json:"id"`
}

func (h *handler) createCheckoutSession(email string, plan customer.Plan, w http.ResponseWriter, r *http.Request) {
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
		SuccessURL: stripe.String("https://woningfinder.nl/start/bedankt"),
		CancelURL:  stripe.String("https://woningfinder.nl/start/geannuleerd"),
	}

	session, err := session.New(params)
	if err != nil {
		errorMsg := fmt.Errorf("error while creating stripe new checkout session")
		h.logger.Sugar().Warnf("%w: %w", errorMsg, err)
		render.Render(w, r, handlerErrors.ServerErrorRenderer(errorMsg))
		return
	}

	// return response
	json.NewEncoder(w).Encode(createCheckoutSessionResponse{
		SessionID: session.ID,
	})
}

// planToLineItems gets the plan price and converts it to a stripe price
// note 1€ is 100 for stripe
func planToLineItems(plan customer.Plan) *stripe.CheckoutSessionLineItemParams {
	switch plan {
	case customer.PlanBasis:
		return &stripe.CheckoutSessionLineItemParams{
			Currency: stripe.String(string(stripe.CurrencyEUR)),
			Name:     stripe.String("Basis"),
			Amount:   stripe.Int64(int64(customer.PlanBasis.Price()) * 100),
			Quantity: stripe.Int64(1),
		}
	case customer.PlanPro:
		return &stripe.CheckoutSessionLineItemParams{
			Currency: stripe.String(string(stripe.CurrencyEUR)),
			Name:     stripe.String("Pro"),
			Amount:   stripe.Int64(int64(customer.PlanPro.Price()) * 100),
			Quantity: stripe.Int64(1),
		}
	}

	return nil
}
