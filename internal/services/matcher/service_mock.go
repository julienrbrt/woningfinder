package matcher

import "github.com/woningfinder/woningfinder/internal/domain/entity"

type serviceMock struct {
	Service
	err error
}

// NewServiceMock mocks the matcher service
func NewServiceMock(err error) Service {
	return &serviceMock{err: err}
}

func (s *serviceMock) MatchOffer(offerList entity.OfferList) error {
	return nil
}
