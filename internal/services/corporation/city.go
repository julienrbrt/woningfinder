package corporation

import (
	"fmt"
	"strings"

	"github.com/woningfinder/woningfinder/internal/domain/entity"
)

func (s *service) AddCities(cities []entity.City, corporation entity.Corporation) error {
	_, err := s.dbClient.Conn().Model(&cities).OnConflict("(name) DO UPDATE").Insert()
	if err != nil {
		return fmt.Errorf("error creating cities: %w", err)
	}

	// add cities relation
	for _, city := range cities {
		city.Name = strings.Title(city.Name)
		if _, err := s.dbClient.Conn().Model(&entity.CorporationCity{CorporationName: corporation.Name, CityName: city.Name}).
			Where("corporation_name = ? and city_name = ?", corporation.Name, city.Name).
			SelectOrInsert(); err != nil {
			return fmt.Errorf("failing creating corporation: %w", err)
		}
	}

	return nil
}

func (s *service) GetCity(name string) (*entity.City, error) {
	var city entity.City
	if err := s.dbClient.Conn().Model(&city).Where("name = ?", name).Select(); err != nil {
		return nil, fmt.Errorf("failing getting city %s: %w", name, err)
	}

	return &city, nil
}

func (s *service) GetCities() (*[]entity.City, error) {
	var cities []entity.City

	if err := s.dbClient.Conn().Model(&cities).Select(); err != nil {
		return nil, fmt.Errorf("failing getting cities: %w", err)
	}

	return &cities, nil
}

func (s *service) DeleteCity(city entity.City) error {
	// TODO to implement
	panic("not implemented")
}
