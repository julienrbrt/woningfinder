package notifications

import (
	"fmt"

	"github.com/woningfinder/woningfinder/internal/auth"

	"github.com/matcornic/hermes/v2"
	"github.com/woningfinder/woningfinder/internal/domain/entity"
	"github.com/woningfinder/woningfinder/internal/services/notifications/templates"
)

func (s *service) SendWelcome(user *entity.User) error {
	_, jwtToken, err := auth.CreateJWTUserToken(s.jwtAuth, user)
	if err != nil {
		return fmt.Errorf("error sending welcome notification: %w", err)
	}

	html, plain, err := welcomeTpl(user, jwtToken)
	if err != nil {
		return fmt.Errorf("error sending welcome notification: %w", err)
	}

	if err := s.emailClient.Send("Welkom bij WoningFinder!", html, plain, user.Email); err != nil {
		return fmt.Errorf("error sending welcome notification: %w", err)
	}

	return nil
}

func welcomeTpl(user *entity.User, jwtToken string) (html, plain string, err error) {
	email := hermes.Email{
		Body: hermes.Body{
			Title: "Welkom bij WoningFinder!",
			Intros: []string{
				"Je zoekopdracht is ingesteld.",
				"Om voor jou te kunnen reageren, hoef je alleen maar in te loggen bij de woningcorporaties waar je wilt reageren.",
				"We hebben voor jou alleen de woningcorporaties geselecteerd die met jouw zoekopdracht matchen.",
			},
			Actions: []hermes.Action{
				{
					Instructions: "We raden je aan bij elke in te loggen zodat je sneller een huis kunt vinden.",
					Button: hermes.Button{
						Color: "#E46948",
						Text:  "Mijn zoekopdracht",
						Link:  fmt.Sprintf("https://woningfinder.nl/mijn-zoekopdracht?jwt=%s", jwtToken),
					},
				},
			},
			Outros: []string{
				"Dan kan je ontspannen, we reageren voor jou. Je hoeft verder niets meer te doen.",
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
