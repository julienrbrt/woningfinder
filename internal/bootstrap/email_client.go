package bootstrap

import (
	"github.com/woningfinder/woningfinder/pkg/config"
	"github.com/woningfinder/woningfinder/pkg/email"
)

// CreateEmailClient creates an email client
func CreateEmailClient() email.Client {
	return email.NewClient(config.MustGetString("APP_NAME"),
		config.MustGetString("EMAIL_ADDRESS"),
		config.MustGetString("EMAIL_PASSWORD"),
		config.MustGetString("EMAIL_SMTP_ADDRESS"),
		config.MustGetInt("EMAIL_SMTP_PORT"),
	)
}
