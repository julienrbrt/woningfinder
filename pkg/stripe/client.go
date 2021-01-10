package stripe

import (
	"net/url"

	"github.com/stripe/stripe-go"
	"github.com/woningfinder/woningfinder/internal/domain/entity"
	"go.uber.org/zap"
)

// APIEndpoint is the Stripe base endpoint
var APIEndpoint = url.URL{Scheme: "https", Host: "api.stripe.com"}

type Client interface {
	CreateCustomer(u *entity.User) error
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
