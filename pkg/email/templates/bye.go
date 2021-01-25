package templates

import (
	"fmt"

	"github.com/matcornic/hermes/v2"
)

// ByeSuject is the subject of the goodbye mail
const ByeSuject = "Hoera je hebt een huis gevonden!"

func (t *templates) ByeTpl() (html, plain string, err error) {
	email := hermes.Email{
		Body: hermes.Body{
			Title: fmt.Sprintf("Hallo %s,", t.user.Name),
			Intros: []string{
				"We hebben gezien dat je een huis hebt gevonden. Van harte gefeliciteerd!",
				"Omdat je ons niet meer nodig hebt, zijn al je gegevens van ons systeem verwijderd (we zijn pricacy-freaks, weet je nog).",
				"Bedankt voor je vertrouwen in ons en geniet van jouw nieuwe woning.",
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
