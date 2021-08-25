package email

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/woningfinder/woningfinder/internal/auth"
	"github.com/woningfinder/woningfinder/internal/customer"
)

func (s *service) SendWelcome(user *customer.User) error {
	_, jwtToken, err := auth.CreateJWTUserToken(s.jwtAuth, user)
	if err != nil {
		return fmt.Errorf("error sending welcome email: %w", err)
	}

	body := &bytes.Buffer{}
	tpl := template.Must(template.ParseFS(emailTemplates, "templates/welcome.html"))
	if err := tpl.Execute(body, fmt.Sprintf("https://woningfinder.nl/mijn-zoekopdracht?jwt=%s", jwtToken)); err != nil {
		return fmt.Errorf("error sending welcome email: %w", err)
	}

	if err := s.emailClient.Send("Welkom bij WoningFinder!", body.String(), user.Email); err != nil {
		return fmt.Errorf("error sending welcome email: %w", err)
	}

	return nil
}
