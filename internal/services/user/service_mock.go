package user

import "github.com/woningfinder/woningfinder/internal/domain/entity"

type serviceMock struct {
	Service
	err error
}

// NewServiceMock mocks the user service
func NewServiceMock(err error) Service {
	return &serviceMock{err: err}
}

func (s *serviceMock) CreateUser(u *entity.User) error {
	return s.err
}

// ExpectedMockGetHousingPreferencesMatchingCorporation is returned when mocking GetHousingPreferencesMatchingCorporation from userServoce
var ExpectedMockGetHousingPreferencesMatchingCorporation = []entity.Corporation{{Name: "De Woonplaats"}}

func (s *serviceMock) GetHousingPreferencesMatchingCorporation(_ *entity.User) ([]entity.Corporation, error) {
	if s.err != nil {
		return nil, s.err
	}

	return ExpectedMockGetHousingPreferencesMatchingCorporation, nil
}

func (s *serviceMock) CreateCorporationCredentials(_ *entity.User, _ entity.CorporationCredentials) error {
	return s.err
}
