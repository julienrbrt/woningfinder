package user

import (
	"time"

	"github.com/woningfinder/woningfinder/internal/corporation"
	"github.com/woningfinder/woningfinder/internal/corporation/city"
	"github.com/woningfinder/woningfinder/internal/customer"
)

type serviceMock struct {
	Service
	err error
}

// NewServiceMock mocks the user service
func NewServiceMock(err error) Service {
	return &serviceMock{err: err}
}

func (s *serviceMock) CreateUser(_ *customer.User) error {
	return s.err
}

func (s *serviceMock) ConfirmUser(email string) (*customer.User, error) {
	return s.GetUser(&customer.User{Email: email})
}

func (s *serviceMock) ConfirmPayment(email string) (*customer.User, error) {
	return s.GetUser(&customer.User{Email: email})
}

func (s *serviceMock) GetUser(search *customer.User) (*customer.User, error) {
	return &customer.User{
		ID:           search.ID,
		Name:         "Test",
		Email:        "test@example.org",
		BirthYear:    1990,
		YearlyIncome: 30000,
		FamilySize:   3,
		Plan: customer.UserPlan{
			CreatedAt:          time.Date(2021, 12, 31, 1, 1, 0, 0, time.UTC),
			Name:               customer.PlanBasis.Name,
			FreeTrialStartedAt: time.Date(2099, 12, 31, 15, 1, 0, 0, time.UTC),
		},
		HousingPreferences: customer.HousingPreferences{
			Type: []corporation.HousingType{
				corporation.HousingTypeHouse,
				corporation.HousingTypeAppartement,
			},
			MaximumPrice:  950,
			NumberBedroom: 1,
			HasElevator:   true,
			City: []city.City{
				{Name: "Enschede"},
			},
		},
	}, s.err
}

func (s *serviceMock) UpdateHousingPreferences(user *customer.User, housingPreferences customer.HousingPreferences) error {
	return s.err
}

func (s *serviceMock) DeleteHousingPreferences(user *customer.User) error {
	return s.err
}

func (s *serviceMock) GetHousingPreferencesMatchingCorporation(_ *customer.User) ([]corporation.Corporation, error) {
	if s.err != nil {
		return nil, s.err
	}

	return []corporation.Corporation{{Name: "De Woonplaats", URL: "https://dewoonplaats.nl"}}, nil
}

func (s *serviceMock) CreateCorporationCredentials(_ uint, _ customer.CorporationCredentials) error {
	return s.err
}

func (s *serviceMock) GetCorporationCredentials(userID uint, corporationName string) (*customer.CorporationCredentials, error) {
	return &customer.CorporationCredentials{}, s.err
}

func (s *serviceMock) CreateWaitingList(w *customer.WaitingList) error {
	return s.err
}
