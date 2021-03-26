package matcher

import (
	"context"

	"github.com/woningfinder/woningfinder/internal/domain/entity"
)

type serviceMock struct {
	Service
	err error
}

// NewServiceMock mocks the matcher service
func NewServiceMock(err error) Service {
	return &serviceMock{err: err}
}

func (s *serviceMock) MatchOffer(_ context.Context, _ entity.OfferList) error {
	return nil
}