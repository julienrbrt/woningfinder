package user

import (
	"fmt"
	"strings"

	"github.com/woningfinder/woningfinder/internal/entity"
)

func (s *service) GetHousingPreferencesMatchingCorporation(u *entity.User) ([]entity.Corporation, error) {
	db := s.dbClient.Conn()

	housingPreferences, err := s.GetHousingPreferences(u)
	if err != nil {
		return nil, fmt.Errorf("error when getting matching corporations: %w", err)
	}

	// get corporation relevant to user housing preferences
	var corporationsMatch []entity.CorporationCity
	if err := db.Model(&corporationsMatch).Where(fmt.Sprintf("city_name IN (%s)", buildCityList(housingPreferences))).DistinctOn("corporation_name").Select(); err != nil {
		return nil, fmt.Errorf("error when getting matching corporations: %w", err)
	}

	var corporations []entity.Corporation
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

// buildCityList extract the cities from the user housing preferences
func buildCityList(housingPreferences entity.HousingPreferences) string {
	var cities []string
	for _, city := range housingPreferences.City {
		cities = append(cities, fmt.Sprintf("'%s'", city.Name))
	}

	return strings.Join(cities, ",")
}
