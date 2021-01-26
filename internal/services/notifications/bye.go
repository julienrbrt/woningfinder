package notifications

import (
	"fmt"

	"github.com/matcornic/hermes/v2"
	"github.com/woningfinder/woningfinder/internal/domain/entity"
	"github.com/woningfinder/woningfinder/internal/services/notifications/templates"
)

func (s *service) SendBye(user *entity.User) error {
	html, plain, err := byeTpl(user)
	if err != nil {
		return fmt.Errorf("error sending bye notification: %w", err)
	}

	if err := s.emailClient.Send("Hoera je hebt een huis gevonden!", html, plain, user.Email); err != nil {
		return fmt.Errorf("error sending bye notification: %w", err)
	}

	return nil
}

func byeTpl(user *entity.User) (html, plain string, err error) {
	email := hermes.Email{
		Body: hermes.Body{
			Title: fmt.Sprintf("Hallo %s,", user.Name),
			Intros: []string{
				"We hebben gezien dat je een huis hebt gevonden. Van harte gefeliciteerd!",
				"Omdat je ons niet meer nodig hebt, zijn al je gegevens van ons systeem verwijderd (we zijn pricacy-freaks, weet je nog).",
				"Bedankt voor je vertrouwen in ons en geniet van jouw nieuwe woning.",
			},
			Signature: "Groetjes",
		},
	}

	// generate html email
	html, err = templates.WoningFinderInfo.GenerateHTML(email)
	if err != nil {
		return "", "", fmt.Errorf("error while building email from template: %w", err)
	}

	// generate plain email
	plain, err = templates.WoningFinderInfo.GeneratePlainText(email)
	if err != nil {
		return "", "", fmt.Errorf("error while building email from template: %w", err)
	}

	return
}
