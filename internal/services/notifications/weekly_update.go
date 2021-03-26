package notifications

import (
	"fmt"

	"github.com/matcornic/hermes/v2"
	"github.com/woningfinder/woningfinder/internal/domain/entity"
	"github.com/woningfinder/woningfinder/internal/services/notifications/templates"
)

func (s *service) SendWeeklyUpdate(user *entity.User, housingMatch []entity.HousingPreferencesMatch) error {
	html, plain, err := weeklyUpdateTpl(user, housingMatch)
	if err != nil {
		return fmt.Errorf("error sending weekly update notification: %w", err)
	}

	if err := s.emailClient.Send("Wekelijkse update", html, plain, user.Email); err != nil {
		return fmt.Errorf("error sending weekly update notification: %w", err)
	}

	return nil
}

func weeklyUpdateTpl(user *entity.User, housingMatch []entity.HousingPreferencesMatch) (html, plain string, err error) {
	var email hermes.Email

	if len(housingMatch) > 0 {
		// list the housing match
		var matchs [][]hermes.Entry
		for _, match := range housingMatch {
			matchs = append(matchs, []hermes.Entry{
				{Key: "Reactie datum", Value: fmt.Sprintf("%d-%d-%d", match.CreatedAt.Day(), match.CreatedAt.Month(), match.CreatedAt.Year())},
				{Key: "Adres", Value: match.HousingAddress},
				{Key: "Woningcorporatie", Value: match.CorporationName},
			})
		}

		introText := "We hebben goed nieuws! In de afgelopen week hebben we op een woning voor jou gereageerd:"
		if len(housingMatch) > 1 {
			introText = fmt.Sprintf("We hebben goed nieuws! In de afgelopen week hebben we op %d woningen voor jou gereageerd:", len(housingMatch))
		}

		email = hermes.Email{
			Body: hermes.Body{
				Title: fmt.Sprintf("Hallo %s,", user.Name),
				Intros: []string{
					introText,
				},
				Table: hermes.Table{
					Data: matchs,
					Columns: hermes.Columns{
						CustomWidth: map[string]string{
							"Woningcorporatie": "20%",
						},
					},
				},
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
				Title: fmt.Sprintf("Hallo %s,", user.Name),
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
