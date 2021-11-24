package corporation

import "github.com/julienrbrt/woningfinder/internal/corporation/city"

type serviceMock struct {
	Service

	expectedMockGetCities []*city.City
	err                   error
}

// NewServiceMock mocks the corporation service
func NewServiceMock(err error) Service {
	return &serviceMock{
		err: err,
		expectedMockGetCities: []*city.City{
			&city.Enschede,
			&city.Hengelo,
		},
	}
}

func (s *serviceMock) GetCities() ([]*city.City, error) {
	if s.err != nil {
		return nil, s.err
	}

	return s.expectedMockGetCities, nil
}
