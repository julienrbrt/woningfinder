package email

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/woningfinder/woningfinder/internal/customer"
)

func (s *service) SendThankYou(user *customer.User) error {
	body := &bytes.Buffer{}
	tpl := template.Must(template.ParseFS(emailTemplates, "templates/thank-you.html"))

	if err := tpl.Execute(body, user.Name); err != nil {
		return fmt.Errorf("error sending thank you email: %w", err)
	}

	if err := s.emailClient.Send("Je WoningFinder zoekopdracht is nu onpeberkt geldig!", body.String(), user.Email); err != nil {
		return fmt.Errorf("error sending thank you email: %w", err)
	}

	return nil
}
