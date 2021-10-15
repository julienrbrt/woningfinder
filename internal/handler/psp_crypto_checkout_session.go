package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/render"
	"github.com/woningfinder/woningfinder/internal/customer"
	handlerErrors "github.com/woningfinder/woningfinder/internal/handler/errors"
	"github.com/woningfinder/woningfinder/pkg/cryptocom"
)

func (h *handler) createCryptoCheckoutSession(email string, plan customer.Plan, w http.ResponseWriter, r *http.Request) {
	params := cryptocom.CryptoCheckoutSession{
		Currency:  "EUR",
		ReturnURL: successURL,
		CancelURL: cancelURL,
		Metadata: cryptocom.CustomerData{
			Email: email,
		},
	}

	params = enrichWithPlan(params, plan)

	response, err := h.cryptoClient.CreatePayment(params)
	if err != nil {
		errorMsg := fmt.Errorf("error while creating crypto.com new checkout session")
		h.logger.Sugar().Errorf("%w: %w", errorMsg, err)
		render.Render(w, r, handlerErrors.ServerErrorRenderer(errorMsg))
		return
	}

	// return response
	json.NewEncoder(w).Encode(subscriptionResponse{
		CryptoPaymentURL: response.PaymentURL,
	})
}

func enrichWithPlan(params cryptocom.CryptoCheckoutSession, plan customer.Plan) cryptocom.CryptoCheckoutSession {
	switch plan {
	case customer.PlanPro:
		params.Description = strings.Title(customer.PlanPro.Name)
		params.Amount = customer.PlanPro.Price * 100 // price in cents
	}

	return params
}
