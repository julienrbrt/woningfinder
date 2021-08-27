package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/stripe/stripe-go"
	"github.com/woningfinder/woningfinder/internal/customer"
)

type cryptoCheckoutSession struct {
	ID             string `json:"id"`
	Amount         int    `json:"amount"`
	AmountRefunded int    `json:"amount_refunded"`
	Created        int    `json:"created"`
	CryptoCurrency string `json:"crypto_currency"`
	CryptoAmount   string `json:"crypto_amount"`
	Currency       string `json:"currency"`
	CustomerID     string `json:"customer_id"`
	PaymentURL     string `json:"payment_url"`
	ReturnURL      string `json:"return_url"`
	CancelURL      string `json:"cancel_url"`
	Description    string `json:"description"`
	LiveMode       bool   `json:"live_mode"`
	Metadata       struct {
		CustomerName string `json:"customer_name"`
	} `json:"metadata"`
	OrderID   string `json:"order_id"`
	Recipient string `json:"recipient"`
	Refunded  bool   `json:"refunded"`
	Status    string `json:"status"`
}

func (h *handler) createCryptoCheckoutSession(email string, plan customer.Plan, w http.ResponseWriter, r *http.Request) {
	params := &cryptoCheckoutSession{
		Currency:  string(stripe.CurrencyEUR),
		ReturnURL: successURL,
		CancelURL: cancelURL,
	}

	switch plan {
	case customer.PlanBasis:
		params.Description = strings.Title(customer.PlanBasis.Name)
		params.Amount = customer.PlanBasis.Price
	case customer.PlanPro:
		params.Description = strings.Title(customer.PlanPro.Name)
		params.Amount = customer.PlanPro.Price
	}

	var response cryptoCheckoutSession
	// make call to crypto.com

	// return response
	json.NewEncoder(w).Encode(paymentProcessorResponse{
		CryptoPaymentURL: response.PaymentURL,
	})
}
