package user

import "github.com/woningfinder/woningfinder/internal/entity"

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

func (s *serviceMock) GetUser(search *entity.User) (*entity.User, error) {
	return &entity.User{
		ID:           search.ID,
		Name:         "Test",
		Email:        "test@example.org",
		BirthYear:    1990,
		YearlyIncome: 30000,
		FamilySize:   3,
		Plan: entity.UserPlan{
			Name: entity.PlanBasis,
		},
		HousingPreferences: entity.HousingPreferences{
			Type: []entity.HousingType{
				entity.HousingTypeHouse,
				entity.HousingTypeAppartement,
			},
			MaximumPrice:  950,
			NumberBedroom: 1,
			HasElevator:   true,
			City: []entity.City{
				{Name: "Enschede"},
			},
		},
	}, s.err
}

func (s *serviceMock) GetHousingPreferencesMatchingCorporation(_ *entity.User) ([]entity.Corporation, error) {
	if s.err != nil {
		return nil, s.err
	}

	return []entity.Corporation{{Name: "De Woonplaats", URL: "https://dewoonplaats.nl"}}, nil
}

func (s *serviceMock) CreateCorporationCredentials(_ uint, _ entity.CorporationCredentials) error {
	return s.err
}

func (s *serviceMock) GetCorporationCredentials(userID uint, corporationName string) (*entity.CorporationCredentials, error) {
	return &entity.CorporationCredentials{}, s.err
}
