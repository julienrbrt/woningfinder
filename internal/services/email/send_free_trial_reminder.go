package email

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/woningfinder/woningfinder/internal/customer"
)

func (s *service) SendFreeTrialReminder(user *customer.User) error {
	plan, err := customer.PlanFromName(user.Plan.Name)
	if err != nil {
		return fmt.Errorf("error sending end free trial reminder email: %w", err)
	}

	data := struct {
		Name, URL string
		PlanPrice int
	}{
		Name:      user.Name,
		URL:       fmt.Sprintf("https://woningfinder.nl/start/voltooien?email=%s", user.Email),
		PlanPrice: plan.Price,
	}

	body := &bytes.Buffer{}
	tpl := template.Must(template.ParseFS(emailTemplates, "templates/free-trial-reminder.html"))
	if err := tpl.Execute(body, data); err != nil {
		return fmt.Errorf("error sending end free trial reminder email: %w", err)
	}

	if err := s.emailClient.Send("Je WoningFinder zoekopdracht", body.String(), user.Email); err != nil {
		return fmt.Errorf("error sending end free trial reminder email: %w", err)
	}

	return nil
}
