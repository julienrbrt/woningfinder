package email

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/julienrbrt/woningfinder/internal/customer"
)

func (s *service) SendBye(user *customer.User) error {
	body := &bytes.Buffer{}
	tpl := template.Must(template.ParseFS(emailTemplates, "templates/bye.html"))

	if err := tpl.Execute(body, user.Name); err != nil {
		return fmt.Errorf("error sending bye email: %w", err)
	}

	if err := s.emailClient.Send("Hoera je hebt een huis gevonden!", body.String(), user.Email); err != nil {
		return fmt.Errorf("error sending bye email: %w", err)
	}

	return nil
}
