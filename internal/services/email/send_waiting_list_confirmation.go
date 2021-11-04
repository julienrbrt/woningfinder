package email

import (
	"bytes"
	"fmt"
	"html/template"
)

func (s *service) SendWaitingListConfirmation(email, city string) error {
	body := &bytes.Buffer{}
	tpl := template.Must(template.ParseFS(emailTemplates, "templates/waiting-list-confirmation.html"))

	if err := tpl.Execute(body, city); err != nil {
		return fmt.Errorf("error sending waiting list confirmation email: %w", err)
	}

	if err := s.emailClient.Send("Je staat op ons wachtlijst!", body.String(), email); err != nil {
		return fmt.Errorf("error sending waiting list confirmation email: %w", err)
	}

	return nil
}
