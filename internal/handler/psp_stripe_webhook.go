package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/go-chi/render"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/webhook"
	"github.com/woningfinder/woningfinder/internal/customer"
	handlerErrors "github.com/woningfinder/woningfinder/internal/handler/errors"
)

const stripeHeader = "Stripe-Signature"

// StripeWebhook is called via the Stripe webhook and confirm that a user has paid
func (h *handler) StripeWebhook(w http.ResponseWriter, r *http.Request) {
	// get request
	const MaxBodyBytes = int64(65536)
	r.Body = http.MaxBytesReader(w, r.Body, MaxBodyBytes)
	payload, err := ioutil.ReadAll(r.Body)
	if err != nil {
		render.Render(w, r, &handlerErrors.ErrorResponse{
			Err:        err,
			StatusCode: http.StatusServiceUnavailable,
			StatusText: "Service Unavailable",
			Message:    fmt.Sprintf("Error reading request body: %v", err)})
		return
	}

	// parse event
	event := stripe.Event{}
	if err := json.Unmarshal(payload, &event); err != nil {
		render.Render(w, r, handlerErrors.ErrorRenderer(fmt.Errorf("failed to parse webhook body json: %w", err)))
		return
	}

	// verify stripe authenticity
	signatureHeader := r.Header.Get(stripeHeader)
	event, err = webhook.ConstructEvent(payload, signatureHeader, h.paymentWebhookSigningKey)
	if err != nil {
		render.Render(w, r, handlerErrors.ErrorRenderer(fmt.Errorf("⚠️ Webhook signature verification failed: %w", err)))
		return
	}

	// check if customer successfully paid
	if event.Type == "payment_intent.succeeded" {
		var paymentIntent stripe.PaymentIntent
		err := json.Unmarshal(event.Data.Raw, &paymentIntent)
		if err != nil {
			render.Render(w, r, handlerErrors.ErrorRenderer(fmt.Errorf("failed to parse webhook json: %w", err)))
			return
		}

		// populate payment - note 1€ is 100 for stripe
		plan := customer.PlanFromPrice(paymentIntent.Amount / 100)
		if plan == (customer.Plan{}) {
			h.logger.Sugar().Warnf("⚠️ Unknown amount %d€ paid by %s: %w", paymentIntent.Amount/100, paymentIntent.ReceiptEmail, err)
			return
		}

		// set payment as proceed
		if err := h.paymentService.ProcessPayment(paymentIntent.ReceiptEmail, plan); err != nil {
			errorMsg := fmt.Errorf("error while processing payment")
			h.logger.Sugar().Warnf("%w: %w", errorMsg, err)
			render.Render(w, r, handlerErrors.ServerErrorRenderer(errorMsg))
			return
		}

		h.logger.Sugar().Infof("🎉🎉🎉 New customer %s paid %d€ 🎉🎉🎉", paymentIntent.ReceiptEmail, paymentIntent.Amount/100)
	}

	// returns 200 by default
}
