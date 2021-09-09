package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/go-chi/render"
	stripeGo "github.com/stripe/stripe-go/v72"
	webhook "github.com/stripe/stripe-go/v72/webhook"
	"github.com/woningfinder/woningfinder/internal/customer"
	handlerErrors "github.com/woningfinder/woningfinder/internal/handler/errors"
	"github.com/woningfinder/woningfinder/pkg/stripe"
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
	event := stripeGo.Event{}
	if err := json.Unmarshal(payload, &event); err != nil {
		render.Render(w, r, handlerErrors.BadRequestErrorRenderer(fmt.Errorf("failed to parse webhook body json: %w", err)))
		return
	}

	// verify stripe authenticity
	signatureHeader := r.Header.Get(stripeHeader)
	event, err = webhook.ConstructEvent(payload, signatureHeader, h.stripeClient.WebhookSigningKey())
	if err != nil {
		render.Render(w, r, handlerErrors.BadRequestErrorRenderer(fmt.Errorf("⚠️ Webhook signature verification failed: %w", err)))
		return
	}

	// check if customer successfully paid
	if event.Type == stripe.PaymentIntentSucceeded {
		var paymentIntent stripeGo.PaymentIntent
		err := json.Unmarshal(event.Data.Raw, &paymentIntent)
		if err != nil {
			render.Render(w, r, handlerErrors.BadRequestErrorRenderer(fmt.Errorf("failed to parse webhook json: %w", err)))
			return
		}

		// populate payment - 1€ is 100 cents
		_, err = customer.PlanFromPrice(paymentIntent.Amount / 100)
		if err != nil {
			h.logger.Sugar().Errorf("⚠️ Unknown amount %d€ paid by %s: %w", paymentIntent.Amount/100, paymentIntent.ReceiptEmail, err)
			return
		}

		// set payment as proceed
		user, err := h.userService.ConfirmPayment(paymentIntent.ReceiptEmail)
		if err != nil {
			errorMsg := fmt.Errorf("error while processing payment")
			h.logger.Sugar().Warnf("%w: %w", errorMsg, err)
			render.Render(w, r, handlerErrors.ServerErrorRenderer(errorMsg))
			return
		}

		// send payment confirmation email
		if err := h.emailService.SendThankYou(user); err != nil {
			h.logger.Sugar().Warn(err)
		}

		h.logger.Sugar().Infof("🎉🎉🎉 New customer %s paid %d€ 🎉🎉🎉", paymentIntent.ReceiptEmail, paymentIntent.Amount/100)
	}

	// returns 200 by default
}
