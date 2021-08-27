package cryptocom

import (
	"net/url"

	"github.com/woningfinder/woningfinder/pkg/networking"
)

// APIEndpoint is the Crypto.com Pay api endpoint
var APIEndpoint = url.URL{Scheme: "https", Host: "pay.crypto.com", Path: "/api/payments"}

// Client for Crypto.com Pay
type Client interface {
	CreatePayment(session CryptoCheckoutSession) (*CryptoCheckoutSession, error)
}

type client struct {
	networkingClient networking.Client
}

// NewClient creates a client for Crypto.com
func NewClient(c networking.Client) Client {
	return &client{
		networkingClient: c,
	}
}
