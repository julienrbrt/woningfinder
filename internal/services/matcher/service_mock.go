package matcher

type serviceMock struct {
	Service
	err error
}

// NewServiceMock mocks the matcher service
func NewServiceMock(err error) Service {
	return &serviceMock{err: err}
}
