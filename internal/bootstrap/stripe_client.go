package bootstrap

import (
	"github.com/julienrbrt/woningfinder/pkg/config"
	"github.com/julienrbrt/woningfinder/pkg/logging"
	"github.com/julienrbrt/woningfinder/pkg/stripe"
)

// CreateSripeClient initialize stripe with the right credentials
func CreateSripeClient(logger *logging.Logger) stripe.Client {
	return stripe.NewClient(logger, config.MustGetString("STRIPE_API_KEY"), config.MustGetString("STRIPE_WEBHOOK_SIGNING_KEY"))
}
