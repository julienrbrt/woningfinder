package email

import "github.com/woningfinder/woningfinder/internal/customer"

type serviceMock struct {
	err error
}

// NewServiceMock mocks the email service
func NewServiceMock(err error) Service {
	return &serviceMock{err: err}
}

func (s *serviceMock) SendActivationEmail(_ *customer.User) error {
	return s.err
}

func (s *serviceMock) SendEmailConfirmationReminder(_ *customer.User) error {
	return s.err
}

func (s *serviceMock) SendThankYou(_ *customer.User) error {
	return s.err
}

func (s *serviceMock) SendPaymentFailed(_ *customer.User) error {
	return s.err
}

func (s *serviceMock) SendLogin(_ *customer.User) error {
	return s.err
}

func (s *serviceMock) SendWeeklyUpdate(_ *customer.User, _ []*customer.HousingPreferencesMatch) error {
	return s.err
}

func (s *serviceMock) SendCorporationCredentialsFirstTimeAdded(_ *customer.User) error {
	return s.err
}

func (s *serviceMock) SendCorporationCredentialsMissing(_ *customer.User) error {
	return s.err
}

func (s *serviceMock) SendCorporationCredentialsError(_ *customer.User, _ string) error {
	return s.err
}

func (s *serviceMock) SendBye(_ *customer.User) error {
	return s.err
}

func (s *serviceMock) ContactFormSubmission(name, email, message string) error {
	return s.err
}
