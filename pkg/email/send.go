package email

import (
	"fmt"
	"net/smtp"

	"github.com/jordan-wright/email"
)

func (c *client) Send(subject, htmlBody, to string) error {
	e := &email.Email{
		To:      []string{to},
		From:    fmt.Sprintf("WoningFinder <%s>", c.from),
		Subject: subject,
		HTML:    []byte(htmlBody),
	}

	if err := e.Send(fmt.Sprintf("%s:%d", c.server, c.port), smtp.PlainAuth("", c.username, c.password, c.server)); err != nil {
		return fmt.Errorf("error while sending mail to %s: %w", to, err)
	}

	return nil
}
