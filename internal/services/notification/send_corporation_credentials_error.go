package notification

import (
	"fmt"

	"github.com/matcornic/hermes/v2"
	"github.com/woningfinder/woningfinder/internal/auth"
	"github.com/woningfinder/woningfinder/internal/customer"
	"github.com/woningfinder/woningfinder/internal/services/notification/templates"
)

func (s *service) SendCorporationCredentialsError(user *customer.User, corporationName string) error {
	_, jwtToken, err := auth.CreateJWTUserToken(s.jwtAuth, user)
	if err != nil {
		return fmt.Errorf("error sending corporation credentials notification: %w", err)
	}

	html, plain, err := corporationCredentialsErrorTpl(user, jwtToken, corporationName)
	if err != nil {
		return fmt.Errorf("error sending corporation credentials notification: %w", err)
	}

	if err := s.emailClient.Send("Er is iets misgegaan met je inloggegevens", html, plain, user.Email); err != nil {
		return fmt.Errorf("error sending corporation credentials notification: %w", err)
	}

	return nil
}

func corporationCredentialsErrorTpl(user *customer.User, jwtToken, corporationName string) (html, plain string, err error) {
	email := hermes.Email{
		Body: hermes.Body{
			Title: fmt.Sprintf("Hoi %s,", user.Name),
			Intros: []string{
				fmt.Sprintf("We hebben geprobeerd om in te loggen bij %s, maar het lijkt erop dat je inloggegevens niet meer kloppen (je hebt waarschijnlijk je wachtwoord veranderd).", corporationName),
				"Ze zijn nu dus verwijderd van ons systeem.",
			},
			Actions: []hermes.Action{
				{
					Instructions: fmt.Sprintf("Als je nog steeds wilt dat we reageren op het %s aanbod log dan in om je inloggegevens voor deze woningcorporatie opnieuw in te stellen.", corporationName),
					Button: hermes.Button{
						Color: "#E46948",
						Text:  "Mijn zoekopdracht",
						Link:  fmt.Sprintf("https://woningfinder.nl/mijn-zoekopdracht?jwt=%s", jwtToken),
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
