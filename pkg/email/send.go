package email

import (
	"fmt"

	"github.com/mattevans/postmark-go"
)

func (c *client) Send(subject, htmlBody, to string) error {
	emailReq := &postmark.Email{
		From:       "WoningFinder <contact@woningfinder.nl>",
		To:         to,
		Subject:    subject,
		HTMLBody:   htmlBody,
		TrackOpens: true,
	}

	if _, response, err := c.postmark.Email.Send(emailReq); err != nil {
		return fmt.Errorf("failed to send email, got response %v: %w", response, err)
	}

	return nil
}
