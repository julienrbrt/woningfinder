package templates

import (
	"fmt"

	"github.com/matcornic/hermes/v2"
)

// WelcomeSubject is the subject of the welcome mail
const WelcomeSubject = "Welkom bij WoningFinder!"

func (t *templates) WelcomeTpl() (html, plain string, err error) {
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
						Text:  "Mijn woningcorporaties",
						Link:  fmt.Sprintf("https://app.woningfinder.com/woningcorporaties?jwt=%s", t.jwtToken),
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
	html, err = t.product.GenerateHTML(email)
	if err != nil {
		return "", "", fmt.Errorf("error while building email from template: %w", err)
	}

	// generate plain email
	plain, err = t.product.GeneratePlainText(email)
	if err != nil {
		return "", "", fmt.Errorf("error while building email from template: %w", err)
	}

	return
}
