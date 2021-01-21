package corporation

import (
	"fmt"

	"github.com/woningfinder/woningfinder/internal/domain/entity"
)

func (s *service) AddCities(cities []entity.City) ([]entity.City, error) {
	_, err := s.dbClient.Conn().Model(&cities).OnConflict("(name) DO NOTHING").Insert()
	if err != nil {
		return nil, fmt.Errorf("error creating cities: %w", err)
	}

	return cities, nil
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
