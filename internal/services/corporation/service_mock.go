package corporation

import "github.com/woningfinder/woningfinder/internal/corporation"

type serviceMock struct {
	Service
	err error
}

// NewServiceMock mocks the corporation service
func NewServiceMock(err error) Service {
	return &serviceMock{err: err}
}

// ExpectedMockGetCities is returned when mocking GetCities from corporationService
var ExpectedMockGetCities = []corporation.City{
	{Name: "Enschede", District: []string{"Roombeek", "Centrum"}},
	{Name: "Hengelo"},
}

func (s *serviceMock) GetCities() ([]corporation.City, error) {
	if s.err != nil {
		return nil, s.err
	}

	return ExpectedMockGetCities, nil
}
