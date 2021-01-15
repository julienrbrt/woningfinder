package bootstrap

import (
	"github.com/woningfinder/woningfinder/pkg/config"
	"github.com/woningfinder/woningfinder/pkg/logging"
	"github.com/woningfinder/woningfinder/pkg/stripe"
)

// CreateSripeClient creates a client for Stripe
func CreateSripeClient(logger *logging.Logger) stripe.Client {
	return stripe.NewClient(logger, config.MustGetString("STRIPE_API_KEY"))
}
