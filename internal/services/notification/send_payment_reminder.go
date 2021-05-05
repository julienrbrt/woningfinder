package notification

import (
	"fmt"

	"github.com/matcornic/hermes/v2"
	"github.com/woningfinder/woningfinder/internal/customer"
	"github.com/woningfinder/woningfinder/internal/services/notification/templates"
)

func (s *service) SendPaymentReminder(user *customer.User) error {
	html, plain, err := paymentReminder(user)
	if err != nil {
		return fmt.Errorf("error sending payment reminder notification: %w", err)
	}

	if err := s.emailClient.Send("Je WoningFinder zoekopdracht", html, plain, user.Email); err != nil {
		return fmt.Errorf("error sending payment reminder notification: %w", err)
	}

	return nil
}

func paymentReminder(user *customer.User) (html, plain string, err error) {
	email := hermes.Email{
		Body: hermes.Body{
			Title: "Je WoningFinder zoekopdracht",
			Intros: []string{
				fmt.Sprintf("Hallo %s,", user.Name),
				"We hebben gezien dat je een zoekopdracht hebt ingesteld maar nog niet hebt voltooid.",
			},
			Actions: []hermes.Action{
				{
					Instructions: "Je kunt dat gemakkelijk hier doen:",
					Button: hermes.Button{
						Color: "#E46948",
						Text:  "Zoekopdracht voltooien",
						Link:  "https://woningfinder.nl/start/geannuleerd",
					},
				},
			},
			Outros: []string{
				"Hulp nodig of vragen? Antwoord deze e-mail, we helpen je graag.",
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
