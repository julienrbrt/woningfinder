package corporation

import "github.com/woningfinder/woningfinder/internal/city"

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
			{Name: "Enschede", District: map[string][]string{"Centrum": nil, "Noord": {"Roombeek"}}},
			{Name: "Hengelo"},
		},
	}
}

func (s *serviceMock) GetCities() ([]*city.City, error) {
	if s.err != nil {
		return nil, s.err
	}

	return s.expectedMockGetCities, nil
}
