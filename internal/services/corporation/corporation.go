package corporation

import (
	"fmt"

	"github.com/woningfinder/woningfinder/internal/domain/entity"
)

// TODO eventually use a prepare function to create it in one query only
func (s *service) CreateOrUpdateCorporation(corp entity.Corporation) error {
	db := s.dbClient.Conn()

	// verify corporation
	if err := corp.IsValid(); err != nil {
		return fmt.Errorf("failing creating corporation %v: %w", corp, err)
	}

	// creates the corporation - on data changes update it
	if _, err := db.Model(&corp).OnConflict("(name) DO UPDATE").Insert(); err != nil {
		return fmt.Errorf("failing creating corporation: %w", err)
	}

	// add corporation selection method
	for _, selection := range corp.SelectionMethod {
		if _, err := db.Model(&entity.CorporationSelectionMethod{CorporationName: corp.Name, SelectionMethod: selection.Method}).
			Where("corporation_name = ? and selection_method = ?", corp.Name, selection.Method).
			SelectOrInsert(); err != nil {
			return fmt.Errorf("failing creating corporation: %w", err)
		}
	}

	// add cities and cities relation
	if err := s.AddCities(corp.Cities, corp); err != nil {
		return fmt.Errorf("failing adding cities to corporation: %w", err)
	}

	return nil
}

func (s *service) GetCorporation(name string) (*entity.Corporation, error) {
	db := s.dbClient.Conn()

	var corp entity.Corporation
	if err := db.Model(&corp).Where("name ILIKE ?", name).Select(); err != nil {
		return nil, fmt.Errorf("failed getting corporation %s: %w", name, err)
	}

	// enriching corporation
	if err := db.Model(&corp).Where("name = ?", corp.Name).Relation("SelectionMethod").Select(); err != nil {
		return nil, fmt.Errorf("failed getting corporation %s selection method: %w", name, err)
	}
	if err := db.Model(&corp).Where("name = ?", corp.Name).Relation("Cities").Select(); err != nil {
		return nil, fmt.Errorf("failed getting corporation %s cities: %w", name, err)
	}

	return &corp, nil
}

func (s *service) DeleteCorporation(corp entity.Corporation) error {
	// TODO to implement
	// Delete all relationships and delete newly unsupported cities
	panic("not implemented")
}
