package email

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/woningfinder/woningfinder/internal/customer"
)

func (s *service) SendWeeklyUpdate(user *customer.User, housingMatch []customer.HousingPreferencesMatch) error {
	body := &bytes.Buffer{}

	if len(housingMatch) == 0 {
		tpl := template.Must(template.ParseFS(emailTemplates, "templates/weekly-update-empty.html"))
		if err := tpl.Execute(body, user.Name); err != nil {
			return fmt.Errorf("error sending weekly update email: %w", err)
		}
	} else {
		data := struct {
			Name        string
			NumberMatch int
			Match       []customer.HousingPreferencesMatch
		}{
			Name:        user.Name,
			NumberMatch: len(housingMatch),
			Match:       housingMatch,
		}

		tpl := template.Must(template.New("weekly-update.html").
			Funcs(template.FuncMap{"mod": func(i, j int) bool { return i%j == 0 }}).
			ParseFS(emailTemplates, "templates/weekly-update.html"))
		if err := tpl.Execute(body, data); err != nil {
			return fmt.Errorf("error sending weekly update email: %w", err)
		}
	}

	if err := s.emailClient.Send("Wekelijkse update", body.String(), user.Email); err != nil {
		return fmt.Errorf("error sending weekly update email: %w", err)
	}

	return nil
}
