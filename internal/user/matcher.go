package user

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/woningfinder/woningfinder/internal/corporation"
	"github.com/woningfinder/woningfinder/internal/database"
)

func (s *service) MatchOffer(offerList corporation.OfferList) error {
	// create housing corporation client
	client, err := s.clientProvider.Get(offerList.Corporation)
	if err != nil {
		return err
	}

	// find users corporation credentials for this offers
	credentials, err := s.findMatchCredentials(offerList)
	if err != nil {
		// no users found, exit silently
		if errors.Is(err, errNoMatchFound) {
			return nil
		}
		return fmt.Errorf("error while matching offer: %w", err)
	}

	// find users with housing preferences matching offer
	users, err := s.listUsersFromCredentials(credentials)
	if err != nil {
		return fmt.Errorf("error while matching offer: %w", err)
	}

	// build credentials map
	var credentialsMap = make(map[int]CorporationCredentials, len(credentials))
	for _, c := range credentials {
		credentialsMap[c.UserID] = c
	}

	for _, user := range users {
		// react concurrently
		go func(user User) {
			//get housing preferences
			housingPreferences, err := s.GetHousingPreferences(&user)
			if err != nil {
				s.logger.Sugar().Errorf("error while getting housing preferences for user %s: %w", user.Email, err)
				return
			}
			user.HousingPreferences = *housingPreferences

			// decrypt housing corporation credentials
			creds := credentialsMap[int(user.ID)]
			creds, err = s.decryptCredentials(creds)
			if err != nil {
				s.logger.Sugar().Errorf("error while decrypting credentials for %s: %w", user.Email, err)
				return
			}

			// login to housing corporation
			if err := client.Login(creds.Login, creds.Password); err != nil {
				s.logger.Sugar().Errorf("failed to login to corporation %s for %s: %w", offerList.Corporation.Name, user.Email, err)
				return
			}

			for _, offer := range offerList.Offer {
				s.logger.Sugar().Infof("checking match of %s for %s...", offer.Housing.Address, user.Email)

				// check if we already check this offer
				uuid := buildReactionUUID(&user, offer)
				if s.hasReacted(uuid) {
					s.logger.Sugar().Debug("has already been checked... skipping.")
					continue
				}

				if user.MatchPreferences(offer) && user.MatchCriteria(offer) {
					// apply
					if err := client.ReactToOffer(offer); err != nil {
						s.logger.Sugar().Errorf("failed to react to %s with user %s: %w", offer.Housing.Address, user.Email, err)
						continue
					}

					// TODO add to queue to send mail
					s.logger.Sugar().Infof("🎉🎉🎉 WoningFinder has successfully reacted to %s on behalf of %s. 🎉🎉🎉\n", offer.Housing.Address, user.Email)
				}

				// save that we've checked the offer for the user
				s.storeReaction(uuid)
			}
		}(user)
	}

	return nil
}

// hasReacted check if a user already reacted to an offer
func (s *service) hasReacted(uuid string) bool {
	_, err := s.redisClient.Get(uuid)
	if err != nil {
		if !errors.Is(err, database.ErrRedisKeyNotFound) {
			s.logger.Sugar().Errorf("error when getting reaction: %w", err)
		}
		// does not have reacted
		return false
	}

	return true
}

// storeReaction saves that an user reacted to an offer
func (s *service) storeReaction(uuid string) {
	if err := s.redisClient.Set(uuid, true); err != nil {
		s.logger.Sugar().Errorf("error when saving reaction to redis: %w", err)
	}
}

func (s *service) findMatchCredentials(offerList corporation.OfferList) ([]CorporationCredentials, error) {
	var credentials []CorporationCredentials
	var query = CorporationCredentials{
		CorporationName: offerList.Corporation.Name,
		CorporationURL:  offerList.Corporation.URL,
	}
	if err := s.dbClient.Conn().Model(&CorporationCredentials{}).Where(query).Find(&credentials).Error; err != nil {
		return nil, fmt.Errorf("error of matchOffer while getting user credentials: %w", err)
	}

	// no users found, exit silently
	if len(credentials) == 0 {
		return nil, errNoMatchFound
	}

	return credentials, nil
}

func (s *service) listUsersFromCredentials(credentials []CorporationCredentials) ([]User, error) {
	// get users id
	var usersID []int
	for _, c := range credentials {
		usersID = append(usersID, c.UserID)
	}

	// query each user
	var users []User
	if err := s.dbClient.Conn().Where("id IN ?", usersID).Find(&users).Error; err != nil {
		return nil, fmt.Errorf("error while getting users: %w", err)
	}

	return users, nil
}

func buildReactionUUID(user *User, offer corporation.Offer) string {
	return base64.StdEncoding.EncodeToString([]byte(user.Email + offer.Housing.Address + offer.SelectionDate.String()))
}

// MatchCriteria verifies that an user match the offer criterias
func (u *User) MatchCriteria(offer corporation.Offer) bool {
	age := time.Now().Year() - u.BirthYear
	// checks if offer age is set and check boundaries
	if offer.MinAge > 0 && ((age < offer.MinAge) || (offer.MaxAge != 0 && age > offer.MaxAge)) {
		return false
	}

	// checks if offer family size is set and check boundaries
	if offer.MinFamilySize > 0 && (u.FamilySize < offer.MinFamilySize || (offer.MaxFamilySize > 0 && u.FamilySize > offer.MaxFamilySize)) {
		return false
	}

	// checks if offer incomes is set and check boundaries
	if offer.MinIncome > 0 && u.YearlyIncome > -1 && (u.YearlyIncome < offer.MinIncome || (offer.MaxIncome > 0 && u.YearlyIncome > offer.MaxIncome)) {
		return false
	}

	return true
}

// MatchPreferences verifies that an offer match the user preferences
func (u *User) MatchPreferences(offer corporation.Offer) bool {
	for _, pref := range u.HousingPreferences {
		// match price
		if offer.Housing.Price >= pref.MaximumPrice {
			continue
		}

		// match house type
		if !matchHouseType(pref, offer.Housing) {
			continue
		}

		// match location
		if !matchCity(pref, offer.Housing) || !matchCityDistrict(pref, offer.Housing) {
			continue
		}

		// match characteristics
		if (pref.NumberBedroom > 0 && pref.NumberBedroom > offer.Housing.NumberBedroom) ||
			(pref.HasBalcony && !offer.Housing.Balcony) ||
			(pref.HasGarden && !offer.Housing.Garden) ||
			(pref.HasElevator && !offer.Housing.Elevator) ||
			(pref.HasHousingAllowance && !offer.Housing.HousingAllowance) ||
			(pref.IsAccessible && !offer.Housing.Accessible) ||
			(pref.HasGarage && !offer.Housing.Garage) ||
			(pref.HasAttic && !offer.Housing.Attic) {
			continue
		}

		return true
	}

	return false
}

func matchHouseType(pref HousingPreferences, housing corporation.Housing) bool {
	if len(pref.Type) == 0 {
		return true
	}

	for _, t := range pref.Type {
		if t.Type == housing.Type.Type {
			return true
		}
	}

	return false
}

func matchCity(pref HousingPreferences, housing corporation.Housing) bool {
	if len(pref.City) == 0 {
		return true
	}

	for _, city := range pref.City {
		// prevent that if actual is an empty, then strings.Contains returns true
		if housing.City.Name != "" && strings.Contains(strings.ToLower(housing.City.Name), strings.ToLower(city.Name)) {
			return true
		}
	}

	return false
}

func matchCityDistrict(pref HousingPreferences, housing corporation.Housing) bool {
	// no preferences so accept everything
	if len(pref.CityDistrict) == 0 {
		return true
	}

	// no district but user has preferences so reject
	if housing.CityDistrict.Name == "" {
		return false
	}

	for _, district := range pref.CityDistrict {
		// prevent that if actual is an empty, then strings.Contains returns true
		if housing.CityDistrict.CityName != "" && strings.Contains(strings.ToLower(housing.CityDistrict.CityName), strings.ToLower(district.CityName)) &&
			housing.CityDistrict.Name != "" && strings.Contains(strings.ToLower(housing.CityDistrict.Name), strings.ToLower(district.Name)) {
			return true
		}
	}

	return false
}
