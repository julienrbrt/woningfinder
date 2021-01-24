package corporation

import "github.com/woningfinder/woningfinder/internal/domain/entity"

type serviceMock struct {
	Service
	err error
}

// NewServiceMock mocks the corporation service
func NewServiceMock(err error) Service {
	return &serviceMock{err: err}
}

// ExpectedMockGetCities is returned when mocking GetCities from corporationService
var ExpectedMockGetCities = []entity.City{{Name: "Enschede"}, {Name: "Hengelo"}}

func (s *serviceMock) GetCities() (*[]entity.City, error) {
	if s.err != nil {
		return nil, s.err
	}

	return &ExpectedMockGetCities, nil
}
