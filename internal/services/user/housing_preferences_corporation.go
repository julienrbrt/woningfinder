package user

import (
	"fmt"
	"strings"

	"github.com/woningfinder/woningfinder/internal/domain/entity"
)

func (s *service) GetHousingPreferencesMatchingCorporation(u *entity.User) ([]entity.Corporation, error) {
	db := s.dbClient.Conn()

	user, err := s.GetUser(u)
	if err != nil {
		return nil, fmt.Errorf("error when getting matching corporations: %w", err)
	}

	// get cities wish
	cities := buildCityList(user.HousingPreferences)
	if len(cities) == 0 { // should not happens
		return nil, fmt.Errorf("error when getting matching corporations: no cities found")
	}

	// inner join between corporation cities and user housing preferences cities
	var corporationCities []entity.CorporationCity
	if err := db.Model(&corporationCities).Where("city_name IN (?)", cities).Select(); err != nil {
		return nil, fmt.Errorf("error when getting matching corporations: %w", err)
	}

	var corporations []entity.Corporation
	for _, corporation := range corporationCities {
		corporations = append(corporations, entity.Corporation{Name: corporation.CorporationName})
	}

	return corporations, nil
}

// buildCityList extract the cities from the user housing preferences
func buildCityList(housingPreferences []entity.HousingPreferences) string {
	citiesMap := make(map[string]entity.City)
	for _, pref := range housingPreferences {
		for _, city := range pref.City {
			citiesMap[city.Name] = city
		}
	}

	var cities []string
	for _, city := range citiesMap {
		cities = append(cities, city.Name)
	}

	return strings.Join(cities, ",")
}
