package email

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/woningfinder/woningfinder/internal/customer"
)

func (s *service) SendPaymentReminder(user *customer.User) error {
	body := &bytes.Buffer{}
	tpl := template.Must(template.ParseFS(emailTemplates, "templates/payment-reminder.html"))
	if err := tpl.Execute(body, user.Name); err != nil {
		return fmt.Errorf("error sending payment reminder email: %w", err)
	}

	if err := s.emailClient.Send("Je WoningFinder zoekopdracht", body.String(), user.Email); err != nil {
		return fmt.Errorf("error sending payment reminder email: %w", err)
	}

	return nil
}
