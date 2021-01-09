package stripe

import (
	"net/url"

	"go.uber.org/zap"

	"github.com/woningfinder/woningfinder/internal/user"

	"github.com/stripe/stripe-go"
)

// APIEndpoint is the Stripe base endpoint
var APIEndpoint = url.URL{Scheme: "https", Host: "api.stripe.com"}

type Client interface {
	CreateCustomer(u *user.User) error
}

type client struct {
	logger *zap.Logger
}

func NewClient(logger *zap.Logger, apiKey string) Client {
	// api key
	stripe.Key = apiKey
	// define stripe default logger
	stripe.DefaultLeveledLogger = logger.Sugar()

	return &client{
		logger: logger,
	}
}
