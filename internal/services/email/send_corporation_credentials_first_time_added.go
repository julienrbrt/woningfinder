package email

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/julienrbrt/woningfinder/internal/customer"
)

func (s *service) SendCorporationCredentialsFirstTimeAdded(user *customer.User) error {
	body := &bytes.Buffer{}
	tpl := template.Must(template.ParseFS(emailTemplates, "templates/corporation-credentials-first-time-added.html"))
	if err := tpl.Execute(body, nil); err != nil {
		return fmt.Errorf("error sending corporation credentials first time added email: %w", err)
	}

	if err := s.emailClient.Send("Dat was het!", body.String(), user.Email); err != nil {
		return fmt.Errorf("error sending corporation credentials first time added email: %w", err)
	}

	return nil
}
