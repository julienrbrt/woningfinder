package email

import (
	"fmt"
	"net/smtp"
	"net/textproto"

	"github.com/jordan-wright/email"
)

func (c *client) Send(subject, html, plain, to string) error {
	e := &email.Email{
		To:      []string{to},
		Bcc:     []string{c.from},
		From:    fmt.Sprintf("%s <%s>", c.name, c.from),
		Subject: subject,
		Text:    []byte(plain),
		HTML:    []byte(html),
		Headers: textproto.MIMEHeader{},
	}

	err := e.Send(fmt.Sprintf("%s:%d", c.server, c.port), smtp.PlainAuth("", c.from, c.password, c.server))
	if err != nil {
		return fmt.Errorf("error while sending mail to %s: %w", to, err)
	}

	return nil
}
