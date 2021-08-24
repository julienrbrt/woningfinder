package email

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/woningfinder/woningfinder/internal/auth"
	"github.com/woningfinder/woningfinder/internal/customer"
)

func (s *service) SendWelcome(user *customer.User) error {
	_, jwtToken, err := auth.CreateJWTUserToken(s.jwtAuth, user)
	if err != nil {
		return fmt.Errorf("error sending welcome email: %w", err)
	}

	body := &bytes.Buffer{}
	tpl := template.Must(template.ParseFS(emailTemplates, "templates/welcome.html"))
	if err := tpl.Execute(body, fmt.Sprintf("https://woningfinder.nl/mijn-zoekopdracht?jwt=%s", jwtToken)); err != nil {
		return fmt.Errorf("error sending welcome email: %w", err)
	}

	if err := s.emailClient.Send("Welkom bij WoningFinder!", body.String(), user.Email); err != nil {
		return fmt.Errorf("error sending welcome email: %w", err)
	}

	return nil
}

func (s *service) SendEmailConfirmationReminder(user *customer.User) error {
	_, jwtToken, err := auth.CreateJWTUserToken(s.jwtAuth, user)
	if err != nil {
		return fmt.Errorf("error sending confirmation email: %w", err)
	}

	body := &bytes.Buffer{}
	tpl := template.Must(template.ParseFS(emailTemplates, "templates/email-confirmation-reminder.html"))
	if err := tpl.Execute(body, fmt.Sprintf("https://woningfinder.nl/mijn-zoekopdracht?jwt=%s", jwtToken)); err != nil {
		return fmt.Errorf("error sending confirmation email: %w", err)
	}

	if err := s.emailClient.Send("Je hebt iets vergeten...", body.String(), user.Email); err != nil {
		return fmt.Errorf("error sending confirmation email: %w", err)
	}

	return nil
}

func (s *service) SendThankYou(user *customer.User) error {
	body := &bytes.Buffer{}
	tpl := template.Must(template.ParseFS(emailTemplates, "templates/thank-you.html"))

	if err := tpl.Execute(body, user.Name); err != nil {
		return fmt.Errorf("error sending thank you email: %w", err)
	}

	if err := s.emailClient.Send("Je WoningFinder zoekopdracht is nu onpeberkt geldig!", body.String(), user.Email); err != nil {
		return fmt.Errorf("error sending thank you email: %w", err)
	}

	return nil
}
