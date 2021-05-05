package notification

import "github.com/woningfinder/woningfinder/internal/customer"

type serviceMock struct {
	Service
	err error
}

// NewServiceMock mocks the notification service
func NewServiceMock(err error) Service {
	return &serviceMock{err: err}
}

func (s *serviceMock) SendLogin(_ *customer.User) error {
	return s.err
}
