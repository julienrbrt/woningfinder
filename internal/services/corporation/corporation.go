package corporation

import (
	"fmt"

	"github.com/woningfinder/woningfinder/internal/corporation"
)

// TODO eventually use a prepare function to create it in one query only
func (s *service) CreateOrUpdateCorporation(corp corporation.Corporation) error {
	db := s.dbClient.Conn()

	// verify corporation
	if err := corp.HasMinimal(); err != nil {
		return fmt.Errorf("failing creating corporation %v: %w", corp, err)
	}

	// creates the corporation - on data changes update it
	if _, err := db.Model(&corp).OnConflict("(name) DO UPDATE").Insert(); err != nil {
		return fmt.Errorf("failing creating corporation: %w", err)
	}

	// add cities and cities relation
	if err := s.AAACities(corp.Cities, corp); err != nil {
		return fmt.Errorf("failing adding cities to corporation: %w", err)
	}

	return nil
}

func (s *service) GetCorporation(name string) (*corporation.Corporation, error) {
	db := s.dbClient.Conn()

	var corp corporation.Corporation
	if err := db.Model(&corp).Where("name = ?", name).Select(); err != nil {
		return nil, fmt.Errorf("failed getting corporation %s: %w", name, err)
	}

	// enriching corporation
	var cities []corporation.CorporationCity
	if err := db.Model(&cities).Where("corporation_name = ?", corp.Name).Select(); err != nil {
		return nil, fmt.Errorf("failed getting corporation %s cities: %w", name, err)
	}

	for _, city := range cities {
		corp.Cities = append(corp.Cities, corporation.City{Name: city.CityName})
	}

	return &corp, nil
}

func (s *service) DeleteCorporation(corp corporation.Corporation) error {
	// TODO to implement
	// Delete all relationships and delete newly unsupported cities
	panic("not implemented")
}
