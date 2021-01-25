package templates

import (
	"fmt"

	"github.com/matcornic/hermes/v2"
)

// CorporationCredentialsErrorSubject is the subject of the corporation credentials mail
const CorporationCredentialsErrorSubject = "Er is iets misgegaan met je inloggegevens"

func (t *templates) CorporationCredentialsErrorTpl(corporationName string) (html, plain string, err error) {
	email := hermes.Email{
		Body: hermes.Body{
			Title: fmt.Sprintf("Hoi %s,", t.user.Name),
			Intros: []string{
				fmt.Sprintf("We hebben geprobeerd om in te loggen bij %s, maar het lijkt erop dat je inloggegevens niet meer kloppen (je hebt waarschijnlijk je wachtwoord veranderd).", corporationName),
				"Ze zijn nu dus verwijderd van ons systeem.",
			},
			Actions: []hermes.Action{
				{
					Instructions: fmt.Sprintf("Als je nog steeds wilt dat we reageren op het %s aanbod log dan in om je inloggegevens voor deze woningcorporatie opnieuw in te stellen.", corporationName),
					Button: hermes.Button{
						Color: "#E46948",
						Text:  "Mijn woningcorporaties",
						Link:  fmt.Sprintf("https://app.woningfinder.com/woningcorporaties?jwt=%s", t.jwtToken),
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
