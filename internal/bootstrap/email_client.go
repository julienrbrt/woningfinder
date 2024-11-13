package bootstrap

import (
	"github.com/julienrbrt/woningfinder/pkg/config"
	"github.com/julienrbrt/woningfinder/pkg/email"
)

// CreateEmailClient creates an email client
func CreateEmailClient() email.Client {
	return email.NewClient(config.MustGetString("EMAIL_ADDRESS"),
		config.MustGetString("EMAIL_USERNAME"),
		config.MustGetString("EMAIL_PASSWORD"),
		config.MustGetString("EMAIL_SMTP_ADDRESS"),
		config.MustGetInt("EMAIL_SMTP_PORT"),
	)
}
