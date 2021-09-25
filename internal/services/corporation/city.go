package corporation

import (
	"fmt"
	"strings"

	"github.com/woningfinder/woningfinder/internal/city"
	"github.com/woningfinder/woningfinder/internal/corporation"
)

// LinkCities permits to creates a city and when given associate that city to a corporation
// Note the corporation will not be check when doing the association, always ensure the corporation exists
func (s *service) LinkCities(cities []city.City, corporations ...corporation.Corporation) error {
	_, err := s.dbClient.Conn().Model(&cities).OnConflict("(name) DO UPDATE").Insert()
	if err != nil {
		return fmt.Errorf("error creating cities: %w", err)
	}

	// add cities relation
	for _, city := range cities {
		city.Name = strings.Title(city.Name)
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

func (s *service) GetCity(name string) (*city.City, error) {
	var c city.City
	if err := s.dbClient.Conn().Model(&c).Where("name ILIKE ?", name).Select(); err != nil {
		return nil, fmt.Errorf("failing getting city %s: %w", name, err)
	}

	// enrich city with suggested city districts
	distritcs, ok := city.SuggestedCityDistrictFromName(c.Name)
	if !ok {
		s.logger.Sugar().Warnf("failed to get city district of %s", name)
	}
	c.District = distritcs

	return &c, nil
}

func (s *service) GetCities() ([]*city.City, error) {
	var cities []*city.City

	if err := s.dbClient.Conn().Model(&cities).Select(); err != nil {
		return nil, fmt.Errorf("failing getting cities: %w", err)
	}

	// enrich city with suggested city districts
	for i, c := range cities {
		districts, ok := city.SuggestedCityDistrictFromName(c.Name)
		if !ok {
			s.logger.Sugar().Warnf("failed to get city district of %s", c.Name)
		}
		cities[i].District = districts
	}

	return cities, nil
}
