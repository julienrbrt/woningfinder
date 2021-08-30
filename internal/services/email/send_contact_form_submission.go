package email

import (
	"bytes"
	"fmt"
	"html/template"
)

func (s *service) ContactFormSubmission(name, email, message string) error {
	body := &bytes.Buffer{}
	tpl := template.Must(template.ParseFS(emailTemplates, "templates/contact-form.html"))

	data := struct {
		Name, Email, Message string
	}{
		Name:    name,
		Email:   email,
		Message: message,
	}

	if err := tpl.Execute(body, data); err != nil {
		return fmt.Errorf("error sending contact form: %w", err)
	}

	if err := s.emailClient.Send("WoningFinder Contact Submission", body.String(), "contact@woningfinder.nl"); err != nil {
		return fmt.Errorf("error sending contact form: %w", err)
	}

	return nil
}
