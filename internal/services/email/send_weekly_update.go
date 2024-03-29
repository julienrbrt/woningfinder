package email

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/julienrbrt/woningfinder/internal/customer"
)

func (s *service) SendWeeklyUpdate(user *customer.User, matches []*customer.HousingPreferencesMatch) error {
	body := &bytes.Buffer{}

	if len(matches) == 0 {
		tpl := template.Must(template.ParseFS(emailTemplates, "templates/weekly-update-empty.html"))
		if err := tpl.Execute(body, user.Name); err != nil {
			return fmt.Errorf("error sending weekly update email: %w", err)
		}
	} else {
		// update picture url
		updatePictureURL(matches)

		data := struct {
			Name        string
			NumberMatch int
			Match       []*customer.HousingPreferencesMatch
		}{
			Name:        user.Name,
			NumberMatch: len(matches),
			Match:       matches,
		}

		tpl := template.Must(template.ParseFS(emailTemplates, "templates/weekly-update.html"))
		if err := tpl.Execute(body, data); err != nil {
			return fmt.Errorf("error sending weekly update email: %w", err)
		}
	}

	if err := s.emailClient.Send("Wekelijkse update", body.String(), user.Email); err != nil {
		return fmt.Errorf("error sending weekly update email: %w", err)
	}

	return nil
}

func updatePictureURL(matches []*customer.HousingPreferencesMatch) {
	for _, match := range matches {
		if match.PictureURL == "" {
			match.PictureURL = "https://woningfinder.nl/img-1.png"
			continue
		}

		match.PictureURL = fmt.Sprintf("https://api.woningfinder.nl/match/%s", match.PictureURL)
	}
}
