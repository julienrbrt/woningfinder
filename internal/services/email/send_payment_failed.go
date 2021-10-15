package email

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/woningfinder/woningfinder/internal/customer"
)

func (s *service) SendPaymentFailed(user *customer.User) error {
	data := struct {
		Name, URL string
		PlanPrice int
	}{
		Name:      user.Name,
		URL:       fmt.Sprintf("https://woningfinder.nl/start/voltooien?email=%s", user.Email),
		PlanPrice: customer.PlanPro.Price,
	}

	body := &bytes.Buffer{}
	tpl := template.Must(template.ParseFS(emailTemplates, "templates/payment-failed.html"))
	if err := tpl.Execute(body, data); err != nil {
		return fmt.Errorf("error sending payment failed email: %w", err)
	}

	if err := s.emailClient.Send("Je WoningFinder zoekopdracht", body.String(), user.Email); err != nil {
		return fmt.Errorf("error sending payment failed email: %w", err)
	}

	return nil
}
