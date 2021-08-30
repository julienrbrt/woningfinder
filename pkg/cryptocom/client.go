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
	VerifyEvent(signatureHeader, eventRaw string) bool
}

type client struct {
	networkingClient networking.Client
	apiKey           string
	webookSigningKey string
}

// NewClient creates a client for Crypto.com
func NewClient(c networking.Client, apiKey, webookSigningKey string) Client {
	return &client{
		networkingClient: c,
		apiKey:           apiKey,
		webookSigningKey: webookSigningKey,
	}
}