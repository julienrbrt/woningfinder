package notification

import (
	"fmt"

	"github.com/matcornic/hermes/v2"
	"github.com/woningfinder/woningfinder/internal/auth"
	"github.com/woningfinder/woningfinder/internal/customer"
	"github.com/woningfinder/woningfinder/internal/services/notification/templates"
)

func (s *service) SendLogin(user *customer.User) error {
	_, jwtToken, err := auth.CreateJWTUserToken(s.jwtAuth, user)
	if err != nil {
		return fmt.Errorf("error sending login notification: %w", err)
	}

	html, plain, err := loginTpl(user, jwtToken)
	if err != nil {
		return fmt.Errorf("error sending login notification: %w", err)
	}

	if err := s.emailClient.Send("Jouw WoningFinder login link", html, plain, user.Email); err != nil {
		return fmt.Errorf("error sending welcome notification: %w", err)
	}

	return nil
}

func loginTpl(user *customer.User, jwtToken string) (html, plain string, err error) {
	email := hermes.Email{
		Body: hermes.Body{
			Title: fmt.Sprintf("Hallo %s,", user.Name),
			Intros: []string{
				"Alstublieft jouw link om in te loggen. De link is 6 uur geldig.",
				"Heb je dit niet gevraagd? Je kun deze e-mail weigeren, jouw account blijft veilig.",
			},
			Actions: []hermes.Action{
				{
					Button: hermes.Button{
						Color: "#E46948",
						Text:  "Mijn zoekopdracht",
						Link:  fmt.Sprintf("https://woningfinder.nl/mijn-zoekopdracht?jwt=%s", jwtToken),
					},
				},
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
