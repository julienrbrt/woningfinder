package email

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/julienrbrt/woningfinder/internal/auth"
	"github.com/julienrbrt/woningfinder/internal/customer"
)

func (s *service) SendLogin(user *customer.User) error {
	_, jwtToken, err := auth.CreateJWTUserToken(s.jwtAuth, user)
	if err != nil {
		return fmt.Errorf("error sending login email: %w", err)
	}

	data := struct {
		Name, URL string
	}{
		Name: user.Name,
		URL:  fmt.Sprintf("https://woningfinder.nl/mijn-zoekopdracht?jwt=%s", jwtToken),
	}

	body := &bytes.Buffer{}
	tpl := template.Must(template.ParseFS(emailTemplates, "templates/login.html"))
	if err := tpl.Execute(body, data); err != nil {
		return fmt.Errorf("error sending login email: %w", err)
	}

	if err := s.emailClient.Send("Jouw WoningFinder login link", body.String(), user.Email); err != nil {
		return fmt.Errorf("error sending login email: %w", err)
	}

	return nil
}
