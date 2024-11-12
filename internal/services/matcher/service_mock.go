package matcher

import (
	"context"

	"github.com/julienrbrt/woningfinder/internal/corporation"
)

type serviceMock struct {
	err error
}

// NewServiceMock mocks the matcher service
func NewServiceMock(err error) Service {
	return &serviceMock{err: err}
}

func (s *serviceMock) Match(_ context.Context, _ corporation.Offers) error {
	return s.err
}
