package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/render"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/checkout/session"
	"github.com/woningfinder/woningfinder/internal/domain/entity"
)

type createCheckoutSessionResponse struct {
	SessionID string `json:"id"`
}

func (h *handler) createCheckoutSession(email string, plan entity.Plan, w http.ResponseWriter, r *http.Request) {
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
		render.Render(w, r, entity.ServerErrorRenderer(errorMsg))
	}

	// return response
	json.NewEncoder(w).Encode(createCheckoutSessionResponse{
		SessionID: session.ID,
	})
}

// planToLineItems gets the plan price and converts it to a stripe price
// note 1€ is 100 for stripe
func planToLineItems(plan entity.Plan) *stripe.CheckoutSessionLineItemParams {
	switch plan {
	case entity.PlanBasis:
		return &stripe.CheckoutSessionLineItemParams{
			Currency: stripe.String(string(stripe.CurrencyEUR)),
			Name:     stripe.String("Basis"),
			Amount:   stripe.Int64(int64(entity.PlanBasis.Price()) * 100),
			Quantity: stripe.Int64(1),
		}
	case entity.PlanPro:
		return &stripe.CheckoutSessionLineItemParams{
			Currency: stripe.String(string(stripe.CurrencyEUR)),
			Name:     stripe.String("Pro"),
			Amount:   stripe.Int64(int64(entity.PlanPro.Price()) * 100),
			Quantity: stripe.Int64(1),
		}
	}

	return nil
}
