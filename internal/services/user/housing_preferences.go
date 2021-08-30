package user

import (
	"errors"
	"fmt"
	"strings"

	pg "github.com/go-pg/pg/v10"
	"github.com/woningfinder/woningfinder/internal/corporation"
	"github.com/woningfinder/woningfinder/internal/corporation/city"
	"github.com/woningfinder/woningfinder/internal/customer"
)

// TODO eventually use a prepare function to create it in one query only
func (s *service) CreateHousingPreferences(user *customer.User, housingPreferences customer.HousingPreferences) error {
	db := s.dbClient.Conn()

	// verify housing preferences
	if err := housingPreferences.HasMinimal(); err != nil {
		return err
	}

	// set and verify housing preferences city
	for i, city := range housingPreferences.City {
		newCity, err := s.corporationService.GetCity(strings.Title(city.Name))
		if err != nil {
			return fmt.Errorf("error the given city %s is not supported: %w", city.Name, err)
		}
		housingPreferences.City[i].Name = newCity.Name
	}

	// assign id
	housingPreferences.UserID = user.ID
	// create housing preferences
	if _, err := db.Model(&housingPreferences).Insert(); err != nil {
		return fmt.Errorf("failing adding housing preferences for user %s: %w", user.Email, err)
	}

	// add housing type relations
	for _, housingType := range housingPreferences.Type {
		if _, err := db.Model(&customer.HousingPreferencesHousingType{HousingPreferencesID: housingPreferences.ID, HousingType: string(housingType)}).
			Where("housing_preferences_id = ? and housing_type = ?", housingPreferences.ID, string(housingType)).
			SelectOrInsert(); err != nil {
			return fmt.Errorf("failing adding housing preferences for user %s: %w", user.Email, err)
		}
	}

	// add cities relation
	for _, city := range housingPreferences.City {
		if _, err := db.Model(&customer.HousingPreferencesCity{HousingPreferencesID: housingPreferences.ID, CityName: city.Name}).
			Where("housing_preferences_id = ? and city_name = ?", housingPreferences.ID, city.Name).
			SelectOrInsert(); err != nil {
			return fmt.Errorf("failing adding housing preferences for user %s: %w", user.Email, err)
		}

		// add cities district
		for district := range city.District {
			if _, err := db.Model(&customer.HousingPreferencesCityDistrict{HousingPreferencesID: housingPreferences.ID, CityName: city.Name, Name: district}).
				Where("housing_preferences_id = ? and city_name = ? and name = ?", housingPreferences.ID, city.Name, district).
				SelectOrInsert(); err != nil {
				return fmt.Errorf("failing adding housing preferences for user %s: %w", user.Email, err)
			}
		}
	}

	return nil
}

func (s *service) GetHousingPreferences(user *customer.User) (customer.HousingPreferences, error) {
	db := s.dbClient.Conn()
	var housingPreferences customer.HousingPreferences

	// get housing preferences
	if err := db.Model(&housingPreferences).Where("user_id = ?", user.ID).Select(); err != nil {
		return housingPreferences, fmt.Errorf("error when getting housing preferences for user %s: %w", user.Email, err)
	}

	// enriching housing preferences

	// add the housing types
	var housingTypes []customer.HousingPreferencesHousingType
	if err := db.Model(&housingTypes).Where("housing_preferences_id = ?", housingPreferences.ID).Select(); err != nil {
		return housingPreferences, fmt.Errorf("failed getting user %s housing preferences type: %w", user.Email, err)
	}

	for _, housingType := range housingTypes {
		housingPreferences.Type = append(housingPreferences.Type, corporation.HousingType(housingType.HousingType))
	}

	// add its city districts
	var cityDistricts []customer.HousingPreferencesCityDistrict
	if err := db.Model(&cityDistricts).Where("housing_preferences_id = ?", housingPreferences.ID).Select(); err != nil {
		return housingPreferences, fmt.Errorf("failed getting user %s housing preferences city districts: %w", user.Email, err)
	}

	// add its city
	var cities []customer.HousingPreferencesCity
	if err := db.Model(&cities).Where("housing_preferences_id = ?", housingPreferences.ID).Select(); err != nil {
		return housingPreferences, fmt.Errorf("failed getting user %s housing preferences cities: %w", user.Email, err)
	}

	for _, c := range cities {
		districts := map[string][]string{}
		for _, district := range cityDistricts {
			if district.CityName == c.CityName {
				// this permits to use same districts map[string][]string that is used for suggestion of districts for the user too
				// we however set the districts only as a key, so only the key is used in the preferences matcher
				districts[district.Name] = nil
			}
		}

		housingPreferences.City = append(housingPreferences.City, city.City{Name: c.CityName, District: districts})

	}

	return housingPreferences, nil
}

// TODO build changelog of the housing preferences and only add the right queries
func (s *service) UpdateHousingPreferences(user *customer.User, housingPreferences customer.HousingPreferences) error {
	if err := housingPreferences.HasMinimal(); err != nil {
		return err
	}

	// delete old housing preferences
	if err := s.DeleteHousingPreferences(user); err != nil {
		return err
	}

	// create new housing preferences
	if err := s.CreateHousingPreferences(user, housingPreferences); err != nil {
		return err
	}

	return nil
}

func (s *service) DeleteHousingPreferences(user *customer.User) error {
	db := s.dbClient.Conn()

	if _, err := db.Model((*customer.HousingPreferencesCity)(nil)).Where("housing_preferences_id = ?", user.HousingPreferences.ID).Delete(); err != nil && !errors.Is(err, pg.ErrNoRows) {
		return fmt.Errorf("failed deleting housing preferences cities for %s: %w", user.Email, err)
	}

	if _, err := db.Model((*customer.HousingPreferencesCityDistrict)(nil)).Where("housing_preferences_id = ?", user.HousingPreferences.ID).Delete(); err != nil && !errors.Is(err, pg.ErrNoRows) {
		return fmt.Errorf("failed deleting housing preferences cities districts for %s: %w", user.Email, err)
	}

	if _, err := db.Model((*customer.HousingPreferencesHousingType)(nil)).Where("housing_preferences_id = ?", user.HousingPreferences.ID).Delete(); err != nil && !errors.Is(err, pg.ErrNoRows) {
		return fmt.Errorf("failed deleting housing preferences housing type for %s: %w", user.Email, err)
	}

	if _, err := db.Model((*customer.HousingPreferences)(nil)).Where("user_id = ?", user.ID).Delete(); err != nil && !errors.Is(err, pg.ErrNoRows) {
		return fmt.Errorf("failed deleting housing preferences cities for %s: %w", user.Email, err)
	}

	if _, err := db.Model((*customer.HousingPreferencesMatch)(nil)).Where("user_id = ?", user.ID).Delete(); err != nil && !errors.Is(err, pg.ErrNoRows) {
		return fmt.Errorf("failed deleting housing preferences match for %s: %w", user.Email, err)
	}

	return nil
}

func (s *service) CreateHousingPreferencesMatch(user *customer.User, offer corporation.Offer, corporationName string) error {
	match := customer.HousingPreferencesMatch{
		UserID:          user.ID,
		HousingAddress:  offer.Housing.Address,
		CorporationName: corporationName,
		OfferURL:        offer.URL,
	}

	if _, err := s.dbClient.Conn().Model(&match).Insert(); err != nil {
		return fmt.Errorf("error when adding housing preferences match: %w", err)
	}

	return nil
}
