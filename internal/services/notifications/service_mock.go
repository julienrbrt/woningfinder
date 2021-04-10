package notifications

type serviceMock struct {
	Service
	err error
}

// NewServiceMock mocks the notifications service
func NewServiceMock(err error) Service {
	return &serviceMock{err: err}
}
