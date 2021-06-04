package notification

import (
	"fmt"

	"github.com/matcornic/hermes/v2"
	"github.com/woningfinder/woningfinder/internal/auth"
	"github.com/woningfinder/woningfinder/internal/customer"
	"github.com/woningfinder/woningfinder/internal/services/notification/templates"
)

func (s *service) SendCorporationCredentialsMissing(user *customer.User) error {
	_, jwtToken, err := auth.CreateJWTUserToken(s.jwtAuth, user)
	if err != nil {
		return fmt.Errorf("error sending corporation credentials missing notification: %w", err)
	}

	html, plain, err := corporationCredentialsMissingTpl(user, jwtToken)
	if err != nil {
		return fmt.Errorf("error sending corporation credentials missing notification: %w", err)
	}

	if err := s.emailClient.Send("Wekelijkse update", html, plain, user.Email); err != nil {
		return fmt.Errorf("error sending corporation credentials notification: %w", err)
	}

	return nil
}

func corporationCredentialsMissingTpl(user *customer.User, jwtToken string) (html, plain string, err error) {
	email := hermes.Email{
		Body: hermes.Body{
			Title: fmt.Sprintf("Hi %s,", user.Name),
			Intros: []string{
				"Normaal gesproken zou je een wekelijkse update krijgen, maar je bent nog niet ingelogd op woningcorporaties websites in je WoningFinder account.",
				"Als je dat niet doet kunnen we helaas niet voor jou reageren.",
			},
			Actions: []hermes.Action{
				{
					Instructions: "Je kunt dat eenvoudig doen via je WoningFinder zoekopdracht.",
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
