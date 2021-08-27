package email

func (s *service) ContactFormSubmission(content string) error {
	return s.emailClient.Send("WoningFinder Contact Submission", content, "contact@woningfinder.nl")
}
