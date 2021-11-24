package email

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/julienrbrt/woningfinder/internal/auth"
	"github.com/julienrbrt/woningfinder/internal/customer"
)

func (s *service) SendCorporationCredentialsError(user *customer.User, corporationName string) error {
	_, jwtToken, err := auth.CreateJWTUserToken(s.jwtAuth, user)
	if err != nil {
		return fmt.Errorf("error sending corporation credentials error email: %w", err)
	}

	data := struct {
		Name, CorporationName, URL string
	}{
		Name:            user.Name,
		CorporationName: corporationName,
		URL:             fmt.Sprintf("https://woningfinder.nl/mijn-zoekopdracht?jwt=%s", jwtToken),
	}

	body := &bytes.Buffer{}
	tpl := template.Must(template.ParseFS(emailTemplates, "templates/corporation-credentials-error.html"))
	if err := tpl.Execute(body, data); err != nil {
		return fmt.Errorf("error sending corporation credentials error email: %w", err)
	}

	if err := s.emailClient.Send("Er is iets misgegaan met je inloggegevens", body.String(), user.Email); err != nil {
		return fmt.Errorf("error sending corporation credentials email: %w", err)
	}

	return nil
}
