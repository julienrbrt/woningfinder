package stripe

import (
	"github.com/stripe/stripe-go"
	"github.com/woningfinder/woningfinder/pkg/logging"
)

// Client for Stripe
// Empty because stripe-go does everything, it is only used as a store of keys
type Client interface {
	WebhookSigningKey() string
}

type client struct {
	apiKey           string
	webookSigningKey string
}

// NewClient creates a client for Stripe
func NewClient(logger *logging.Logger, apiKey, webookSigningKey string) Client {
	// set stripe api key
	stripe.Key = apiKey

	// define stripe default logger
	stripe.DefaultLeveledLogger = logger.Sugar()

	return &client{
		apiKey:           apiKey,
		webookSigningKey: webookSigningKey,
	}
}
