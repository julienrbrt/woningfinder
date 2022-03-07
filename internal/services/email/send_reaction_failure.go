package email

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/julienrbrt/woningfinder/internal/auth"
	"github.com/julienrbrt/woningfinder/internal/corporation"
	"github.com/julienrbrt/woningfinder/internal/customer"
)

func (s *service) SendReactionFailure(user *customer.User, corporationName string, offers []corporation.Offer) error {
	_, jwtToken, err := auth.CreateJWTUserToken(s.jwtAuth, user)
	if err != nil {
		return fmt.Errorf("error sending confirmation email: %w", err)
	}

	data := struct {
		Name            string
		CorporationName string
		URL             string
		Offers          []corporation.Offer
	}{
		Name:            user.Name,
		URL:             fmt.Sprintf("https://woningfinder.nl/mijn-zoekopdracht?jwt=%s", jwtToken),
		CorporationName: corporationName,
		Offers:          offers,
	}

	body := &bytes.Buffer{}
	tpl := template.Must(template.ParseFS(emailTemplates, "templates/reaction-failure.html"))
	if err := tpl.Execute(body, data); err != nil {
		return fmt.Errorf("error sending reaction failure email: %w", err)
	}

	if err := s.emailClient.Send(fmt.Sprintf("Er is iets misgegaan met automatisch reageren bij %s", corporationName), body.String(), user.Email); err != nil {
		return fmt.Errorf("error sending reaction failure email: %w", err)
	}

	return nil
}
