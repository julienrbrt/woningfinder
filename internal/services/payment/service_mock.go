package payment

type serviceMock struct {
	Service
	err error
}

// NewServiceMock mocks the payment service
func NewServiceMock(err error) Service {
	return &serviceMock{err: err}
}
