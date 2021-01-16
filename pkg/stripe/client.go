package stripe

import (
	"net/url"

	"github.com/stripe/stripe-go"
	"github.com/woningfinder/woningfinder/internal/domain/entity"
	"github.com/woningfinder/woningfinder/pkg/logging"
)

// APIEndpoint is the Stripe base endpoint
var APIEndpoint = url.URL{Scheme: "https", Host: "api.stripe.com"}

type Client interface {
	CreateCustomer(u *entity.User) error
}

type client struct {
	logger *logging.Logger
}

func NewClient(logger *logging.Logger, apiKey string) Client {
	// api key
	stripe.Key = apiKey
	// define stripe default logger
	stripe.DefaultLeveledLogger = logger.Sugar()

	return &client{
		logger: logger,
	}
}
