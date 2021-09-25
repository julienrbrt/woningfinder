package user

import (
	"errors"
	"fmt"
	"strings"

	pg "github.com/go-pg/pg/v10"
	"github.com/woningfinder/woningfinder/internal/city"
	"github.com/woningfinder/woningfinder/internal/corporation"
	"github.com/woningfinder/woningfinder/internal/customer"
)

// TODO eventually use a prepare function to create it in one query only
func (s *service) CreateHousingPreferences(userID uint, preferences *customer.HousingPreferences) error {
	db := s.dbClient.Conn()

	// verify housing preferences
	if err := preferences.HasMinimal(); err != nil {
		return err
	}

	// set and verify housing preferences city
	for i, city := range preferences.City {
		newCity, err := s.corporationService.GetCity(strings.Title(city.Name))
		if err != nil {
			return fmt.Errorf("error the given city %s is not supported: %w", city.Name, err)
		}
		preferences.City[i].Name = newCity.Name
	}

	// assign id
	preferences.UserID = userID
	// create housing preferences
	if _, err := db.Model(preferences).Insert(); err != nil {
		return fmt.Errorf("failing adding housing preferences for userID %d: %w", userID, err)
	}

	// add housing type relations
	for _, housingType := range preferences.Type {
		if _, err := db.Model(&customer.HousingPreferencesHousingType{UserID: userID, HousingType: string(housingType)}).
			Where("user_id = ? and housing_type = ?", userID, string(housingType)).
			SelectOrInsert(); err != nil {
			return fmt.Errorf("failing adding housing preferences for userID %d: %w", userID, err)
		}
	}

	// add cities relation
	for _, city := range preferences.City {
		if _, err := db.Model(&customer.HousingPreferencesCity{UserID: userID, CityName: city.Name}).
			Where("user_id = ? and city_name = ?", userID, city.Name).
			SelectOrInsert(); err != nil {
			return fmt.Errorf("failing adding housing preferences for userID %d: %w", userID, err)
		}

		// add cities district
		for district := range city.District {
			if _, err := db.Model(&customer.HousingPreferencesCityDistrict{UserID: userID, CityName: city.Name, Name: district}).
				Where("user_id = ? and city_name = ? and name = ?", userID, city.Name, district).
				SelectOrInsert(); err != nil {
				return fmt.Errorf("failing adding housing preferences for userID %d: %w", userID, err)
			}
		}
	}

	return nil
}

// TODO refractor in one request
func (s *service) GetHousingPreferences(userID uint) (customer.HousingPreferences, error) {
	db := s.dbClient.Conn()
	var housingPreferences customer.HousingPreferences

	// get housing preferences
	if err := db.Model(&housingPreferences).Where("user_id = ?", userID).Select(); err != nil {
		return housingPreferences, fmt.Errorf("error when getting housing preferences for userID %d: %w", userID, err)
	}

	// enriching housing preferences

	// add the housing types
	var housingTypes []customer.HousingPreferencesHousingType
	if err := db.Model(&housingTypes).Where("user_id = ?", userID).Select(); err != nil {
		return housingPreferences, fmt.Errorf("failed getting userID %d housing preferences type: %w", userID, err)
	}

	for _, housingType := range housingTypes {
		housingPreferences.Type = append(housingPreferences.Type, corporation.HousingType(housingType.HousingType))
	}

	// add its city districts
	var cityDistricts []customer.HousingPreferencesCityDistrict
	if err := db.Model(&cityDistricts).Where("user_id = ?", userID).Select(); err != nil {
		return housingPreferences, fmt.Errorf("failed getting userID %d housing preferences city districts: %w", userID, err)
	}

	// add its city
	var cities []customer.HousingPreferencesCity
	if err := db.Model(&cities).Where("user_id = ?", userID).Select(); err != nil {
		return housingPreferences, fmt.Errorf("failed getting userID %d housing preferences cities: %w", userID, err)
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

// UpdateHousingPreferences updates the housing preferences of a given user
func (s *service) UpdateHousingPreferences(userID uint, preferences *customer.HousingPreferences) error {
	if userID == 0 {
		return errors.New("userID missing")
	}

	if err := preferences.HasMinimal(); err != nil {
		return err
	}

	// delete old housing preferences
	if err := s.DeleteHousingPreferences(userID); err != nil {
		return err
	}

	// create new housing preferences
	if err := s.CreateHousingPreferences(userID, preferences); err != nil {
		return err
	}

	return nil
}

func (s *service) DeleteHousingPreferences(userID uint) error {
	db := s.dbClient.Conn()

	if _, err := db.Model((*customer.HousingPreferencesCity)(nil)).Where("user_id = ?", userID).Delete(); err != nil && !errors.Is(err, pg.ErrNoRows) {
		return fmt.Errorf("failed deleting housing preferences cities for userID %d: %w", userID, err)
	}

	if _, err := db.Model((*customer.HousingPreferencesCityDistrict)(nil)).Where("user_id = ?", userID).Delete(); err != nil && !errors.Is(err, pg.ErrNoRows) {
		return fmt.Errorf("failed deleting housing preferences cities districts for userID %d: %w", userID, err)
	}

	if _, err := db.Model((*customer.HousingPreferencesHousingType)(nil)).Where("user_id = ?", userID).Delete(); err != nil && !errors.Is(err, pg.ErrNoRows) {
		return fmt.Errorf("failed deleting housing preferences housing type for userID %d: %w", userID, err)
	}

	if _, err := db.Model((*customer.HousingPreferences)(nil)).Where("user_id = ?", userID).Delete(); err != nil && !errors.Is(err, pg.ErrNoRows) {
		return fmt.Errorf("failed deleting housing preferences cities for userID %d: %w", userID, err)
	}

	return nil
}

func (s *service) CreateHousingPreferencesMatch(userID uint, offer corporation.Offer, corporationName string) error {
	if _, err := s.dbClient.Conn().Model(&customer.HousingPreferencesMatch{
		UserID:          userID,
		HousingAddress:  offer.Housing.Address,
		CorporationName: corporationName,
		OfferURL:        offer.URL,
	}).Insert(); err != nil {
		return fmt.Errorf("error when adding housing preferences match: %w", err)
	}

	return nil
}
