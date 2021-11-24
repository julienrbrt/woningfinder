package user

import (
	"time"

	"github.com/julienrbrt/woningfinder/internal/corporation"
	"github.com/julienrbrt/woningfinder/internal/corporation/city"
	"github.com/julienrbrt/woningfinder/internal/customer"
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

func (s *serviceMock) ConfirmUser(_ string) error {
	return s.err
}

func (s *serviceMock) GetUser(email string) (*customer.User, error) {
	return &customer.User{
		ID:           42,
		Name:         "Test",
		Email:        email,
		BirthYear:    1990,
		YearlyIncome: 30000,
		FamilySize:   3,
		Plan: customer.UserPlan{
			CreatedAt:        time.Date(2021, 12, 31, 1, 1, 0, 0, time.UTC),
			Name:             "test-ugly-woningfinder-plan",
			StripeCustomerID: "cus_KQoZm6zke6gelu",
			ActivatedAt:      time.Date(2099, 12, 31, 15, 1, 0, 0, time.UTC),
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

func (s *serviceMock) DeleteUser(_ string) error {
	return s.err
}

func (s *serviceMock) UpdateHousingPreferences(_ uint, _ *customer.HousingPreferences) error {
	return s.err
}

func (s *serviceMock) DeleteHousingPreferences(_ uint) error {
	return s.err
}

func (s *serviceMock) GetHousingPreferencesMatchingCorporation(_ uint) ([]*corporation.Corporation, error) {
	if s.err != nil {
		return nil, s.err
	}

	return []*corporation.Corporation{{Name: "De Woonplaats", URL: "https://dewoonplaats.nl"}}, nil
}

func (s *serviceMock) CreateCorporationCredentials(_ uint, _ *customer.CorporationCredentials) error {
	return s.err
}

func (s *serviceMock) HasCorporationCredentials(_ uint) (bool, error) {
	return false, s.err
}

func (s *serviceMock) GetCorporationCredentials(_ uint, _ string) (*customer.CorporationCredentials, error) {
	return &customer.CorporationCredentials{}, s.err
}

func (s *serviceMock) CreateWaitingList(_ *customer.WaitingList) error {
	return s.err
}

func (s *serviceMock) ConfirmSubscription(_ string) error {
	return s.err
}

func (s *serviceMock) SetStripeCustomerID(_ *customer.User, _ string) error {
	return s.err
}

func (s *serviceMock) UpdateSubscriptionStatus(_ string, _ bool) error {
	return s.err
}

func (s *serviceMock) UpdateUser(_ *customer.User) error {
	return s.err
}
