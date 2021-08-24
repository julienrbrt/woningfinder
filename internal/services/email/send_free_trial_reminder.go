package email

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/woningfinder/woningfinder/internal/customer"
)

func (s *service) SendFreeTrialReminder(user *customer.User) error {
	body := &bytes.Buffer{}
	tpl := template.Must(template.ParseFS(emailTemplates, "templates/free-trial-reminder.html"))
	if err := tpl.Execute(body, user.Name); err != nil {
		return fmt.Errorf("error sending end free trial reminder email: %w", err)
	}

	if err := s.emailClient.Send("Je WoningFinder zoekopdracht", body.String(), user.Email); err != nil {
		return fmt.Errorf("error sending end free trial reminder email: %w", err)
	}

	return nil
}
