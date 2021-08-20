package user

import (
	"fmt"

	pg "github.com/go-pg/pg/v10"
	"github.com/woningfinder/woningfinder/internal/corporation"
	"github.com/woningfinder/woningfinder/internal/customer"
)

func (s *service) GetHousingPreferencesMatchingCorporation(u *customer.User) ([]corporation.Corporation, error) {
	housingPreferences, err := s.GetHousingPreferences(u)
	if err != nil {
		return nil, fmt.Errorf("error when getting matching corporations: %w", err)
	}

	// get corporation relevant to user housing preferences
	var corporationsMatch []corporation.CorporationCity
	if err := s.dbClient.Conn().Model(&corporationsMatch).Where("city_name IN (?)", pg.In(cityList(housingPreferences))).DistinctOn("corporation_name").Select(); err != nil {
		return nil, fmt.Errorf("error when getting matching corporations: %w", err)
	}

	var corporations []corporation.Corporation
	for _, c := range corporationsMatch {
		// enriching corporation
		corporation, err := s.corporationService.GetCorporation(c.CorporationName)
		if err != nil {
			return nil, fmt.Errorf("error failing enriching matching corporations: %w", err)

		}

		corporations = append(corporations, *corporation)
	}

	return corporations, nil
}

// cityList extract the cities from the user housing preferences
func cityList(housingPreferences customer.HousingPreferences) []string {
	var cities []string
	for _, city := range housingPreferences.City {
		cities = append(cities, city.Name)
	}

	return cities
}
