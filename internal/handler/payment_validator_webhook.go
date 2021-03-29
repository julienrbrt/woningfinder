package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/go-chi/render"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/webhook"
	"github.com/woningfinder/woningfinder/internal/domain/entity"
	handlerEntity "github.com/woningfinder/woningfinder/internal/handler/entity"
)

const stripeHeader = "Stripe-Signature"

// ProcessPayment is called via the Stripe webhook and confirm that a user has paid
func (h *handler) ProcessPayment(w http.ResponseWriter, r *http.Request) {
	// get request
	const MaxBodyBytes = int64(65536)
	r.Body = http.MaxBytesReader(w, r.Body, MaxBodyBytes)
	payload, err := ioutil.ReadAll(r.Body)
	if err != nil {
		render.Render(w, r, &handlerEntity.ErrorResponse{
			Err:        err,
			StatusCode: http.StatusServiceUnavailable,
			StatusText: "Service Unavailable",
			Message:    fmt.Sprintf("Error reading request body: %v", err)})
		return
	}

	// parse event
	event := stripe.Event{}
	if err := json.Unmarshal(payload, &event); err != nil {
		render.Render(w, r, handlerEntity.ErrorRenderer(fmt.Errorf("failed to parse webhook body json: %w", err)))
		return
	}

	// verify stripe authenticity
	signatureHeader := r.Header.Get(stripeHeader)
	event, err = webhook.ConstructEvent(payload, signatureHeader, h.paymentWebhookSigningKey)
	if err != nil {
		render.Render(w, r, handlerEntity.ErrorRenderer(fmt.Errorf("⚠️ Webhook signature verification failed: %w", err)))
		return
	}

	// check if customer successfully paid
	if event.Type == "payment_intent.succeeded" {
		var paymentIntent stripe.PaymentIntent
		err := json.Unmarshal(event.Data.Raw, &paymentIntent)
		if err != nil {
			render.Render(w, r, handlerEntity.ErrorRenderer(fmt.Errorf("failed to parse webhook json: %w", err)))
			return
		}

		// populate payment
		plan, err := priceToPlan(paymentIntent.Amount)
		if err != nil {
			h.logger.Sugar().Warnf("⚠️ Unknown amount %d€ paid by %s: %w", paymentIntent.Amount/100, paymentIntent.ReceiptEmail, err)
			return
		}

		// set payment as proceed
		if err := h.paymentService.ProcessPayment(&entity.PaymentData{UserEmail: paymentIntent.ReceiptEmail, Plan: plan}); err != nil {
			render.Render(w, r, handlerEntity.ServerErrorRenderer(fmt.Errorf("error while processing payment: %w", err)))
			return
		}

		h.logger.Sugar().Infof("🎉🎉🎉 New client %s paid %d€ 🎉🎉🎉", paymentIntent.ReceiptEmail, paymentIntent.Amount/100)
	}

	// returns 200 by default
}

// priceToPlan gets the stripe price and converts it to a plan
// note 1€ is 100 for stripe
func priceToPlan(amount int64) (entity.Plan, error) {
	switch amount {
	case 3500:
		return entity.PlanBasis, nil
	case 7500:
		return entity.PlanPro, nil
	}

	return "", errors.New("error plan does not exists")
}
