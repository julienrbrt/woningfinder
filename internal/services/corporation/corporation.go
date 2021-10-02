package corporation

import (
	"fmt"

	"github.com/woningfinder/woningfinder/internal/corporation"
	"github.com/woningfinder/woningfinder/internal/corporation/city"
)

// CreateOrUpdateCorporation creates or update a housing corporation
func (s *service) CreateOrUpdateCorporation(corp corporation.Corporation) error {
	// verify corporation
	if err := corp.HasMinimal(); err != nil {
		return fmt.Errorf("failing creating corporation %v: %w", corp, err)
	}

	// creates the corporation - on data changes update it
	if _, err := s.dbClient.Conn().Model(&corp).OnConflict("(name) DO UPDATE").Insert(); err != nil {
		return fmt.Errorf("failing creating corporation: %w", err)
	}

	// add cities and cities relation
	if err := s.LinkCities(corp.Cities, corp); err != nil {
		return fmt.Errorf("failing adding cities to corporation: %w", err)
	}

	return nil
}

func (s *service) GetCorporation(name string) (*corporation.Corporation, error) {
	var corp corporation.Corporation
	if err := s.dbClient.Conn().Model(&corp).Where("name ILIKE ?", name).Select(); err != nil {
		return nil, fmt.Errorf("failed getting corporation %s: %w", name, err)
	}

	// enriching corporation
	var cities []corporation.CorporationCity
	if err := s.dbClient.Conn().Model(&cities).Where("corporation_name = ?", corp.Name).Select(); err != nil {
		return nil, fmt.Errorf("failed getting corporation %s cities: %w", name, err)
	}

	for _, c := range cities {
		corp.Cities = append(corp.Cities, city.City{Name: c.CityName})
	}

	return &corp, nil
}
