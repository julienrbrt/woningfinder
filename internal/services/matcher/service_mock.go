package matcher

import (
	"context"

	"github.com/woningfinder/woningfinder/internal/corporation"
	"github.com/woningfinder/woningfinder/internal/corporation/connector"
)

type serviceMock struct {
	err error
}

// NewServiceMock mocks the matcher service
func NewServiceMock(err error) Service {
	return &serviceMock{err: err}
}

func (s *serviceMock) MatchOffer(_ context.Context, _ corporation.Offers) error {
	return s.err
}

func (s *serviceMock) PushOffers(client connector.Client, corp corporation.Corporation) error {
	return s.err
}

func (s *serviceMock) SubscribeOffers(ch chan<- corporation.Offers) error {
	return s.err
}
