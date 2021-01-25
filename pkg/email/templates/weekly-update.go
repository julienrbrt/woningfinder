package templates

import (
	"fmt"

	"github.com/matcornic/hermes/v2"
	"github.com/woningfinder/woningfinder/internal/domain/entity"
)

// WeeklyUpdateSubject is the subject of the weekly update mail
const WeeklyUpdateSubject = "Wekelijkse update"

func (t *templates) WeeklyUpdateTpl(housingMatch []entity.HousingPreferencesMatch) (html, plain string, err error) {
	var email hermes.Email

	if len(housingMatch) > 0 {
		email = hermes.Email{
			Body: hermes.Body{
				Title: fmt.Sprintf("Hallo %s,", t.user.Name),
				Intros: []string{
					fmt.Sprintf("We hebben goed nieuws! In de afgelopen week hebben we op %d woning(en) gereageerd:", len(housingMatch)),
				},
				// TODO show houses
				Outros: []string{
					"Voor meer informatie, kun je altijd kijken op de website van de woningcorporaties waar we hebben gereageerd.",
					"Je huis staat tussen jouw reacties.",
					"We hopen dat je word gekozen voor een van deze woningen!",
					"Hulp nodig of vragen? Antwoord deze e-mail, we helpen je graag.",
				},
				Signature: "Groetjes",
			},
		}
	} else {
		email = hermes.Email{
			Body: hermes.Body{
				Title: fmt.Sprintf("Hallo %s,", t.user.Name),
				Intros: []string{
					"We hebben elke dag gekeken, maar hebben deze week niets voor jou kunnen vinden.",
					"Maak je geen zorgen, we blijven zoeken!",
				},
				Outros: []string{
					"Hulp nodig of vragen? Antwoord deze e-mail, we helpen je graag.",
					"Tot volgende week.",
				},
				Signature: "Groetjes",
			},
		}
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
