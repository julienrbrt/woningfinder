package bootstrap

import (
	"github.com/woningfinder/woningfinder/pkg/config"
	"github.com/woningfinder/woningfinder/pkg/logging"
	"github.com/woningfinder/woningfinder/pkg/stripe"
)

// CreateSripeClient initialize stripe with the right credentials
func CreateSripeClient(logger *logging.Logger) stripe.Client {
	return stripe.NewClient(logger, config.MustGetString("STRIPE_API_KEY"), config.MustGetString("STRIPE_WEBHOOK_SIGNING_KEY"))
}
