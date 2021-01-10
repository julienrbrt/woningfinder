package bootstrap

import (
	"github.com/woningfinder/woningfinder/pkg/config"
	"github.com/woningfinder/woningfinder/pkg/stripe"
	"go.uber.org/zap"
)

// CreateSripeClient creates a client for Stripe
func CreateSripeClient(logger *zap.Logger) stripe.Client {
	return stripe.NewClient(logger, config.MustGetString("STRIPE_API_KEY"))
}
