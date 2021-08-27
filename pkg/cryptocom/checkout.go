package cryptocom

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/woningfinder/woningfinder/pkg/networking"
)

type CryptoCheckoutSession struct {
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

func (c *client) CreatePayment(session CryptoCheckoutSession) (*CryptoCheckoutSession, error) {
	body, err := json.Marshal(session)
	if err != nil {
		return nil, fmt.Errorf("error while marshaling %v: %w", session, err)
	}

	request := networking.Request{
		Method: http.MethodPost,
		Body:   bytes.NewBuffer(body),
	}

	resp, err := c.networkingClient.Send(&request)
	if err != nil {
		return nil, err
	}

	var response CryptoCheckoutSession
	err = resp.ReadJSONBody(&response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
