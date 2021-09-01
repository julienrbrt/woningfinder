package user

import (
	"fmt"

	"github.com/woningfinder/woningfinder/internal/corporation"
	"github.com/woningfinder/woningfinder/internal/customer"
)

func (s *service) GetHousingPreferencesMatchingCorporation(userID uint) ([]*corporation.Corporation, error) {
	matchingCities := s.dbClient.Conn().
		Model((*customer.HousingPreferencesCity)(nil)).
		Where("user_id = ?", userID).
		ColumnExpr("city_name")

	// get corporation relevant to user housing preferences
	var corporationsMatch []corporation.CorporationCity
	if err := s.dbClient.Conn().
		Model(&corporationsMatch).
		Where("city_name IN (?)", matchingCities).
		DistinctOn("corporation_name").
		Select(); err != nil {
		return nil, fmt.Errorf("error when getting matching corporations: %w", err)
	}

	var corporations []*corporation.Corporation
	for _, c := range corporationsMatch {
		// enriching corporation
		corporation, err := s.corporationService.GetCorporation(c.CorporationName)
		if err != nil {
			return nil, fmt.Errorf("error failing enriching matching corporations: %w", err)

		}

		corporations = append(corporations, corporation)
	}

	return corporations, nil
}
