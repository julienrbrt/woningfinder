package corporation

import (
	"fmt"
	"strings"

	"github.com/woningfinder/woningfinder/internal/corporation"
)

// LinkCities permits to creates a city and when given associate that city to a corporation
// Note the corporation will not be check when doing the association, always ensure the corporation exists
func (s *service) LinkCities(cities []corporation.City, corporations ...corporation.Corporation) error {
	_, err := s.dbClient.Conn().Model(&cities).OnConflict("(name) DO UPDATE").Insert()
	if err != nil {
		return fmt.Errorf("error creating cities: %w", err)
	}

	// add cities relation
	for _, city := range cities {
		city.Name = strings.Title(city.Name)
		// add cities district to city
		for _, district := range city.District {
			if _, err := s.dbClient.Conn().Model(&corporation.CityDistrict{CityName: city.Name, Name: district.Name}).
				Where("city_name = ? and name = ?", city.Name, district.Name).
				SelectOrInsert(); err != nil {
				return fmt.Errorf("failing adding district to city: %w", err)
			}
		}

		// associate city to corporations
		for _, corp := range corporations {
			if _, err := s.dbClient.Conn().Model(&corporation.CorporationCity{CorporationName: corp.Name, CityName: city.Name}).
				Where("corporation_name = ? and city_name = ?", corp.Name, city.Name).
				SelectOrInsert(); err != nil {
				return fmt.Errorf("failing associating city to corporation: %w", err)
			}
		}
	}

	return nil
}

func (s *service) GetCity(name string) (corporation.City, error) {
	var city corporation.City
	if err := s.dbClient.Conn().Model(&city).Where("name = ?", name).Select(); err != nil {
		return corporation.City{}, fmt.Errorf("failing getting city %s: %w", name, err)
	}

	// enrich city with district
	var districts []corporation.CityDistrict
	if err := s.dbClient.Conn().Model(&districts).Where("city_name = ?", city.Name).Select(); err != nil {
		return corporation.City{}, fmt.Errorf("failing getting city district: %w", err)
	}
	city.District = districts

	return city, nil
}

func (s *service) GetCities() ([]corporation.City, error) {
	var cities []corporation.City

	if err := s.dbClient.Conn().Model(&cities).Select(); err != nil {
		return nil, fmt.Errorf("failing getting cities: %w", err)
	}

	// enrich cities with district
	for i, city := range cities {
		var districts []corporation.CityDistrict
		if err := s.dbClient.Conn().Model(&districts).Where("city_name = ?", city.Name).Select(); err != nil {
			return nil, fmt.Errorf("failing getting city district: %w", err)
		}

		cities[i].District = districts
	}

	return cities, nil
}
