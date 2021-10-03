package email

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/woningfinder/woningfinder/internal/auth"
	"github.com/woningfinder/woningfinder/internal/customer"
)

func (s *service) SendActivationEmail(user *customer.User) error {
	_, jwtToken, err := auth.CreateJWTUserToken(s.jwtAuth, user)
	if err != nil {
		return fmt.Errorf("error sending activation email: %w", err)
	}

	data := struct {
		Name, URL string
	}{
		Name: user.Name,
		URL:  fmt.Sprintf("https://woningfinder.nl/mijn-zoekopdracht?jwt=%s", jwtToken),
	}

	body := &bytes.Buffer{}
	tpl := template.Must(template.ParseFS(emailTemplates, "templates/activation.html"))
	if err := tpl.Execute(body, data); err != nil {
		return fmt.Errorf("error sending activation email: %w", err)
	}

	if err := s.emailClient.Send("Activeer je WoningFinder account!", body.String(), user.Email); err != nil {
		return fmt.Errorf("error sending activation email: %w", err)
	}

	return nil
}
