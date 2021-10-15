package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/go-chi/render"
	"github.com/woningfinder/woningfinder/internal/customer"
	handlerErrors "github.com/woningfinder/woningfinder/internal/handler/errors"
	"github.com/woningfinder/woningfinder/pkg/cryptocom"
)

const cryptoHeader = "Pay-Signature"

// CryptoWebhook is called via the Crypto.com webhook and confirm that a user is subscribed
func (h *handler) CryptoWebhook(w http.ResponseWriter, r *http.Request) {
	// get request
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
	event := cryptocom.WebhookEvent{}
	if err := json.Unmarshal(payload, &event); err != nil {
		render.Render(w, r, handlerErrors.BadRequestErrorRenderer(fmt.Errorf("failed to parse webhook body json: %w", err)))
		return
	}

	// verify stripe authenticity
	signatureHeader := r.Header.Get(cryptoHeader)
	if signatureHeader == "" || !h.cryptoClient.VerifyEvent(signatureHeader, string(payload)) {
		render.Render(w, r, handlerErrors.BadRequestErrorRenderer(errors.New("⚠️ Webhook signature verification failed")))
		return
	}

	switch event.Type {
	case cryptocom.PaymentCaptured:
		// check payment - 1€ is 100 cents
		if _, err = customer.PlanFromPrice(int64(event.Data.Object.Amount / 100)); err != nil {
			h.logger.Sugar().Errorf("⚠️ Unknown amount %d€ paid by %s: %w", event.Data.Object.Amount/100, event.Data.Object.Metadata.Email, err)
			return
		}

		// confirm subscription has payment went through
		user, err := h.userService.ConfirmSubscription(event.Data.Object.Metadata.Email)
		if err != nil {
			errorMsg := fmt.Errorf("error while processing payment")
			h.logger.Sugar().Errorf("%w: %w", errorMsg, err)
			render.Render(w, r, handlerErrors.ServerErrorRenderer(errorMsg))
			return
		}

		// send payment confirmation email
		if err := h.emailService.SendThankYou(user); err != nil {
			h.logger.Sugar().Error(err)
		}

		h.logger.Sugar().Infof("🎉🎉🎉 New customer %s subscribed 🎉🎉🎉", event.Data.Object.Metadata.Email)
	}

	// returns 200 by default
}
