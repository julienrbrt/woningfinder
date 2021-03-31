package user

import (
	"errors"
	"fmt"
	"strings"

	"github.com/woningfinder/woningfinder/internal/domain/entity"
)

// TODO eventually use a prepare function to create it in one query only
func (s *service) CreateHousingPreferences(u *entity.User, preferences []entity.HousingPreferences) error {
	db := s.dbClient.Conn()

	// verify if user allowed to create a housing preferences
	if len(u.HousingPreferences) > u.Plan.Name.MaxHousingPreferences() {
		return fmt.Errorf("error cannot create more than %d housing preferences in plan %s: got %d", u.Plan.Name.MaxHousingPreferences(), u.Plan.Name, len(u.HousingPreferences))
	}

	for _, housingPreferences := range preferences {
		// verify housing preferences
		if len(housingPreferences.Type) == 0 {
			return errors.New("error housing preferences invalid: housing type missing")
		}

		if len(housingPreferences.City) == 0 {
			return errors.New("error housing preferences invalid: cities missing")
		}

		// set and verify housing preferences city
		for i, city := range housingPreferences.City {
			newCity, err := s.corporationService.GetCity(strings.Title(city.Name))
			if err != nil {
				return fmt.Errorf("error the given city %s is not supported: %w", city.Name, err)
			}
			housingPreferences.City[i].Name = newCity.Name
		}

		// assign ids
		housingPreferences.UserID = u.ID
		// create housing preferences
		if _, err := db.Model(&housingPreferences).Insert(); err != nil {
			return fmt.Errorf("failing adding housing preferences for user %s: %w", u.Email, err)
		}

		// add housing type relations
		for _, housingType := range housingPreferences.Type {
			if _, err := db.Model(&entity.HousingPreferencesHousingType{HousingPreferencesID: housingPreferences.ID, HousingType: string(housingType)}).
				Where("housing_preferences_id = ? and housing_type = ?", housingPreferences.ID, string(housingType)).
				SelectOrInsert(); err != nil {
				return fmt.Errorf("failing adding housing preferences for user %s: %w", u.Email, err)
			}
		}

		// add cities relation
		for _, city := range housingPreferences.City {
			if _, err := db.Model(&entity.HousingPreferencesCity{HousingPreferencesID: housingPreferences.ID, CityName: city.Name}).
				Where("housing_preferences_id = ? and city_name = ?", housingPreferences.ID, city.Name).
				SelectOrInsert(); err != nil {
				return fmt.Errorf("failing adding housing preferences for user %s: %w", u.Email, err)
			}

			// add cities district
			for _, district := range city.District {
				if _, err := db.Model(&entity.HousingPreferencesCityDistrict{HousingPreferencesID: housingPreferences.ID, CityName: city.Name, Name: district}).
					Where("housing_preferences_id = ? and city_name = ? and name = ?", housingPreferences.ID, city.Name, district).
					SelectOrInsert(); err != nil {
					return fmt.Errorf("failing adding housing preferences for user %s: %w", u.Email, err)
				}
			}
		}
	}

	return nil
}

func (s *service) GetHousingPreferences(u *entity.User) ([]entity.HousingPreferences, error) {
	db := s.dbClient.Conn()

	// get housing preferences
	if err := db.Model(u).Where("id = ?", u.ID).Relation("HousingPreferences").Select(); err != nil {
		return nil, fmt.Errorf("error when getting housing preferences for user %s: %w", u.Email, err)
	}

	// enrich housing preferences

	// add the types
	if err := db.Model(&u.HousingPreferences).Where("user_id = ?", u.ID).Relation("Type").Select(); err != nil {
		return nil, fmt.Errorf("failed getting user %s housing preferences type: %w", u.Email, err)
	}

	// add its city
	if err := db.Model(&u.HousingPreferences).Where("user_id = ?", u.ID).Relation("City").Select(); err != nil {
		return nil, fmt.Errorf("error when getting housing preferences cities for user %s: %w", u.Email, err)
	}

	// add its city districts
	if err := db.Model(&u.HousingPreferences).Where("user_id = ?", u.ID).Relation("CityDistrict").Select(); err != nil {
		return nil, fmt.Errorf("error when getting housing preferences city districts for user %s: %w", u.Email, err)
	}

	return u.HousingPreferences, nil
}

func (s *service) DeleteHousingPreferences(u *entity.User) error {
	// TODO to implement
	// delete housing preferences
	// delete all relations
	panic("not implemented")
}

func (s *service) CreateHousingPreferencesMatch(u *entity.User, offer entity.Offer, corporationName string) error {
	match := entity.HousingPreferencesMatch{
		UserID:          u.ID,
		HousingAddress:  offer.Housing.Address,
		CorporationName: corporationName,
		OfferURL:        offer.URL,
	}

	if _, err := s.dbClient.Conn().Model(&match).Insert(); err != nil {
		return fmt.Errorf("error when adding housing preferences match: %w", err)
	}

	return nil
}
