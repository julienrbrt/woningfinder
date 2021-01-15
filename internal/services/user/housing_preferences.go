package user

import (
	"fmt"

	"github.com/woningfinder/woningfinder/internal/domain/entity"
)

func (s *service) CreateHousingPreferences(u *entity.User, pref []entity.HousingPreferences) error {
	if !u.Tier.Name.AllowMultipleHousingPreferences() && len(pref) > 1 {
		return fmt.Errorf("error cannot create more than one housing preferences in plan %s", u.Tier.Name)
	}

	for _, housingPreferences := range pref {
		if len(housingPreferences.Type) == 0 {
			return fmt.Errorf("error at least 1 housing type is required to setup a housing preferences")
		}

		// set housing preferences (city and district)
		for i, city := range housingPreferences.City {
			newCity, err := s.corporationService.GetCity(city.Name)
			if err != nil {
				return fmt.Errorf("error the given city %s is not supported", city.Name)
			}

			housingPreferences.City[i].Name = newCity.Name

			for _, district := range city.District {
				district.CityName = newCity.Name
				housingPreferences.CityDistrict = append(housingPreferences.CityDistrict, district)
			}

		}

		// create or replace housing preferences
		if err := s.dbClient.Conn().Model(&u).Association("HousingPreferences").Replace(&housingPreferences); err != nil {
			return fmt.Errorf("error when creating/updating housing preferences: %w", err)
		}
	}

	return nil
}

func (s *service) GetHousingPreferences(u *entity.User) (*[]entity.HousingPreferences, error) {
	var pref []entity.HousingPreferences

	// get housing preferences
	if err := s.dbClient.Conn().Model(u).Association("HousingPreferences").Find(&pref); err != nil {
		return nil, fmt.Errorf("error when getting housing preferences for user %s: %w", u.Email, err)
	}

	for _, housingPreferences := range pref {
		// add the types
		if err := s.dbClient.Conn().Model(housingPreferences).Association("Type").Find(&housingPreferences.Type); err != nil {
			return nil, fmt.Errorf("error when getting housing preferences type for user %s: %w", u.Email, err)
		}

		// add its city
		if err := s.dbClient.Conn().Model(housingPreferences).Association("City").Find(&housingPreferences.City); err != nil {
			return nil, fmt.Errorf("error when getting housing preferences cities for user %s: %w", u.Email, err)
		}

		// add its city districts
		if err := s.dbClient.Conn().Model(housingPreferences).Association("CityDistrict").Find(&housingPreferences.CityDistrict); err != nil {
			return nil, fmt.Errorf("error when getting housing preferences city districts for user %s: %w", u.Email, err)
		}
	}

	return &pref, nil
}

func (s *service) DeleteHousingPreferences(u *entity.User) error {
	pref, err := s.GetHousingPreferences(u)
	if err != nil {
		return fmt.Errorf("error when deleting housing preferences: %w", err)
	}

	return s.dbClient.Conn().
		Unscoped().
		// Select(clause.Associations).
		Delete(pref).Error
}
