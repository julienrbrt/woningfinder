package email

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/julienrbrt/woningfinder/internal/corporation"
	"github.com/julienrbrt/woningfinder/internal/customer"
)

func (s *service) SendReactionFailure(user *customer.User, corporationName string, offer corporation.Offer) error {
	body := &bytes.Buffer{}
	tpl := template.Must(template.ParseFS(emailTemplates, "templates/reaction-failure.html"))

	data := struct {
		Name            string
		CorporationName string
		OfferURL        string
		Address         string
	}{
		Name:            user.Name,
		CorporationName: corporationName,
		OfferURL:        offer.URL,
		Address:         offer.Housing.Address,
	}

	if err := tpl.Execute(body, data); err != nil {
		return fmt.Errorf("error sending reaction failure email: %w", err)
	}

	if err := s.emailClient.Send("Er is iets misgegaan met een reactie", body.String(), user.Email); err != nil {
		return fmt.Errorf("error sending reaction failure email: %w", err)
	}

	return nil
}
