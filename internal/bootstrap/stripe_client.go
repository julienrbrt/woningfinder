package bootstrap

import (
	"github.com/stripe/stripe-go"
	"github.com/woningfinder/woningfinder/pkg/config"
	"github.com/woningfinder/woningfinder/pkg/logging"
)

// CreateSripeClient initialize stripe with the right credentials
func CreateSripeClient(logger *logging.Logger) {
	stripe.Key = config.MustGetString("STRIPE_API_KEY")
	// define stripe default logger
	stripe.DefaultLeveledLogger = logger.Sugar()
}
