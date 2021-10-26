package bootstrap

import (
	"github.com/woningfinder/woningfinder/pkg/config"
	"github.com/woningfinder/woningfinder/pkg/email"
)

// CreateEmailClient creates an email client
func CreateEmailClient() email.Client {
	return email.NewClient(config.MustGetString("POSTMARK_SERVER_KEY"))
}
