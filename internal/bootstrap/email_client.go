package bootstrap

import (
	"github.com/julienrbrt/woningfinder/pkg/config"
	"github.com/julienrbrt/woningfinder/pkg/email"
)

// CreateEmailClient creates an email client
func CreateEmailClient() email.Client {
	return email.NewClient(config.MustGetString("POSTMARK_SERVER_KEY"))
}
