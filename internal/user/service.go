package user

import (
	"encoding/base64"
	"errors"
	"fmt"
	"log"

	"github.com/go-redis/redis"
	"github.com/woningfinder/woningfinder/pkg/aes"

	"gorm.io/gorm/clause"

	"github.com/woningfinder/woningfinder/internal/corporation"
	"gorm.io/gorm"
)

// Service permits to handle the persistence of an user
type Service interface {
	CreateUser(u *User) (*User, error)
	GetUser(email string) (*User, error)
	DeleteUser(u *User) error

	CreateHousingPreferences(u *User, housingPreferences HousingPreferences) error
	GetHousingPreferences(u *User) (*HousingPreferences, error)
	DeleteHousingPreferences(u *User) error

	CreateCorporationCredentials(u *User, credentials CorporationCredentials) error
	GetCorporationCredentials(u *User, corporation corporation.Corporation) (*CorporationCredentials, error)
	DeleteCorporationCredentials(u *User, corporation corporation.Corporation) error

	MatchOffer(offers corporation.OfferList) error
	HasReacted(uuid string) bool
	SaveReaction(uuid string)
}

type userService struct {
	db                 *gorm.DB
	rdb                *redis.Client
	aesSecret          string
	clientProvider     corporation.ClientProvider
	corporationService corporation.Service
}

func NewService(db *gorm.DB, rdb *redis.Client, aesSecret string, clientProvider corporation.ClientProvider, corporationService corporation.Service) Service {
	return &userService{
		db:                 db,
		rdb:                rdb,
		aesSecret:          aesSecret,
		clientProvider:     clientProvider,
		corporationService: corporationService,
	}
}

func (s *userService) CreateUser(u *User) (*User, error) {
	if u.Email == "" {
		return nil, fmt.Errorf("email is required for creating user")
	}

	if u.YearlyIncome < -1 {
		return nil, fmt.Errorf("yearly income must be greater than 0, or set to -1 to not be used")
	}

	if len(u.HousingPreferences.Type) == 0 {
		return nil, fmt.Errorf("housing preferences is required for creating user")
	}

	_, err := s.GetUser(u.Email)
	if err == nil {
		return nil, fmt.Errorf("error user %s already exists", u.Email)
	}

	if err := s.db.Create(&u).Error; err != nil {
		return nil, err
	}

	if err := s.CreateHousingPreferences(u, u.HousingPreferences); err != nil {
		return nil, fmt.Errorf("error when creating user %s: %w", u.Email, err)
	}

	return u, nil
}

func (s *userService) GetUser(email string) (*User, error) {
	var u User
	s.db.Where(&User{Email: email}).First(&u)

	if u.ID == 0 {
		return nil, fmt.Errorf("no user found with the email: %s", email)
	}

	return &u, nil
}

func (s *userService) DeleteUser(u *User) error {
	// delete all corporations credentials
	if err := s.db.Unscoped().Select(clause.Associations).
		Where(&CorporationCredentials{UserID: int(u.ID)}).
		Delete(&CorporationCredentials{}).Error; err != nil {
		return fmt.Errorf("failed to delete corporation credentials associated to this user: %w", err)
	}

	// delete housing preferences
	err := s.DeleteHousingPreferences(u)
	if err != nil {
		return fmt.Errorf("failed deleting housing preferences from user: %w", err)
	}

	return s.db.Unscoped().Select(clause.Associations).Delete(u).Error
}

func (s *userService) CreateHousingPreferences(u *User, housingPreferences HousingPreferences) error {
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
	if err := s.db.Model(&u).Association("HousingPreferences").Replace(&housingPreferences); err != nil {
		return fmt.Errorf("error when creating/updating housing preferences: %w", err)
	}

	return nil
}

func (s *userService) GetHousingPreferences(u *User) (*HousingPreferences, error) {
	var housingPreferences HousingPreferences

	// get housing preferences
	if err := s.db.Model(u).Association("HousingPreferences").Find(&housingPreferences); err != nil {
		return nil, fmt.Errorf("error when getting housing preferences for user %s: %w", u.Email, err)
	}

	// add its city
	if err := s.db.Model(housingPreferences).Association("City").Find(&housingPreferences.City); err != nil {
		return nil, fmt.Errorf("error when getting housing preferences cities for user %s: %w", u.Email, err)
	}

	// add the types
	if err := s.db.Model(housingPreferences).Association("Type").Find(&housingPreferences.Type); err != nil {
		return nil, fmt.Errorf("error when getting housing preferences type for user %s: %w", u.Email, err)
	}

	return &housingPreferences, nil
}

func (s *userService) DeleteHousingPreferences(u *User) error {
	pref, err := s.GetHousingPreferences(u)
	if err != nil {
		return fmt.Errorf("error when deleting housing preferences: %w", err)
	}

	return s.db.
		Unscoped().
		Select(clause.Associations).
		Delete(pref).Error
}

func (s *userService) CreateCorporationCredentials(u *User, credentials CorporationCredentials) error {
	if credentials.Corporation.Name == "" || credentials.Login == "" || credentials.Password == "" {
		return fmt.Errorf("error login or password cannot be empty when adding credentials")
	}

	// check credentials validity
	client, err := s.clientProvider.Get(credentials.Corporation)
	if err != nil {
		return err
	}
	if err := client.Login(credentials.Login, credentials.Password); err != nil {
		return fmt.Errorf("error when authenticating to %s with given credentials: %w", credentials.Corporation.Name, err)
	}

	// encrypt credentials
	credentials.Login, err = aes.Encrypt(credentials.Login, s.aesSecret)
	if err != nil {
		return fmt.Errorf("error when encrypting corporation credentials: %w", err)
	}

	credentials.Password, err = aes.Encrypt(credentials.Password, s.aesSecret)
	if err != nil {
		return fmt.Errorf("error when encrypting corporation credentials: %w", err)
	}

	// check if already existing
	fetchCredentials, err := s.GetCorporationCredentials(u, credentials.Corporation)
	if err != nil && !errors.Is(err, errCorporationCredentialsNotFound) {
		return fmt.Errorf("error when checking if credentials already exists: %w", err)
	}

	// store credentials
	if err != nil { // store unexisting credentials
		if err := s.db.Model(u).Association("CorporationCredentials").Append(&credentials); err != nil {
			return fmt.Errorf("error when creating corporation credentials: %w", err)
		}
	} else { // update existing credentials
		if err := s.db.Model(&fetchCredentials).Updates(&credentials).Error; err != nil {
			return fmt.Errorf("error when updating corporation credentials: %w", err)
		}
	}

	return nil
}

func (s *userService) GetCorporationCredentials(u *User, corporation corporation.Corporation) (*CorporationCredentials, error) {
	query := CorporationCredentials{
		UserID:          int(u.ID),
		CorporationName: corporation.Name,
		CorporationURL:  corporation.URL,
	}

	// get corporation credentials
	var credentials CorporationCredentials
	if err := s.db.Where(query).Find(&credentials).Error; err != nil {
		return nil, fmt.Errorf("error when getting corporation credentials for user %s: %w", u.Email, err)
	}

	if credentials.Login == "" || credentials.Password == "" {
		return nil, errCorporationCredentialsNotFound
	}

	// decrypt credentials
	var err error
	credentials.Login, err = aes.Decrypt(credentials.Login, s.aesSecret)
	if err != nil {
		return nil, fmt.Errorf("error when decrypting corporation credentials: %w", err)
	}

	credentials.Password, err = aes.Decrypt(credentials.Password, s.aesSecret)
	if err != nil {
		return nil, fmt.Errorf("error when decrypting corporation credentials: %w", err)
	}

	return &credentials, nil
}

func (s *userService) DeleteCorporationCredentials(u *User, corporation corporation.Corporation) error {
	credentials, err := s.GetCorporationCredentials(u, corporation)
	if err != nil {
		return fmt.Errorf("error when deleting corporation credentials: %w", err)
	}

	// delete permanently
	credentials.Login = ""
	credentials.Password = ""
	if err = s.db.Unscoped().Delete(credentials).Error; err != nil {
		return fmt.Errorf("error when deleting corporation credentials: %w", err)
	}

	return nil
}

func (s *userService) decryptCredentials(credentials CorporationCredentials) (CorporationCredentials, error) {
	// decrypt credentials
	var err error
	credentials.Login, err = aes.Decrypt(credentials.Login, s.aesSecret)
	if err != nil {
		return CorporationCredentials{}, fmt.Errorf("error when decrypting corporation credentials: %w", err)
	}

	credentials.Password, err = aes.Decrypt(credentials.Password, s.aesSecret)
	if err != nil {
		return CorporationCredentials{}, fmt.Errorf("error when decrypting corporation credentials: %w", err)
	}

	return credentials, nil
}

func (s *userService) MatchOffer(offerList corporation.OfferList) error {
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
				log.Printf("error while getting housing preferences for user %s: %v\n", user.Email, err)
				return
			}
			user.HousingPreferences = *housingPreferences

			// decrypt housing corporation credentials
			creds := credentialsMap[int(user.ID)]
			creds, err = s.decryptCredentials(creds)
			if err != nil {
				log.Printf("error while decrypting credentials for %s: %v\n", user.Email, err)
				return
			}

			// login to housing corporation
			if err := client.Login(creds.Login, creds.Password); err != nil {
				log.Printf("failed to login to corporation %s for %s: %v\n", offerList.Corporation.Name, user.Email, err)
				return
			}

			for _, offer := range offerList.Offer {
				log.Printf("checking match of %s for %s...", offer.Housing.Address, user.Email)

				// check if we already check this offer
				uuid := buildReactionUUID(&user, offer)
				if s.HasReacted(uuid) {
					log.Println("has already been checked... skipping.")
					continue
				}

				if user.MatchPreferences(offer) && user.MatchCriteria(offer) {
					// apply
					if err := client.ReactToOffer(offer); err != nil {
						log.Printf("failed to react to %s with user %s: %v\n", offer.Housing.Address, user.Email, err)
						continue
					}

					// TODO add to queue to send mail
					log.Printf("🎉🎉🎉 WoningFinder has successfully reacted to %s on behalf of %s. 🎉🎉🎉\n", offer.Housing.Address, user.Email)
				}

				// save that WoningFinder checks that offer for this user
				s.SaveReaction(uuid)
			}
		}(user)
	}

	return nil
}

func (s *userService) findMatchCredentials(offerList corporation.OfferList) ([]CorporationCredentials, error) {
	var credentials []CorporationCredentials
	var query = CorporationCredentials{
		CorporationName: offerList.Corporation.Name,
		CorporationURL:  offerList.Corporation.URL,
	}
	if err := s.db.Model(&CorporationCredentials{}).Where(query).Find(&credentials).Error; err != nil {
		return nil, fmt.Errorf("error of matchOffer while getting user credentials: %w", err)
	}

	// no users found, exit silently
	if len(credentials) == 0 {
		return nil, errNoMatchFound
	}

	return credentials, nil
}

func (s *userService) listUsersFromCredentials(credentials []CorporationCredentials) ([]User, error) {
	// get users id
	var usersID []int
	for _, c := range credentials {
		usersID = append(usersID, c.UserID)
	}

	// query each user
	var users []User
	if err := s.db.Where("id IN ?", usersID).Find(&users).Error; err != nil {
		return nil, fmt.Errorf("error while getting users: %w", err)
	}

	return users, nil
}

// HasReacted check if a user already reacted to an offer
func (s *userService) HasReacted(uuid string) bool {
	_, err := s.rdb.Get(uuid).Result()
	if err != nil {
		// if error is different that key does not exists
		if err != redis.Nil {
			log.Printf("error when getting reaction from redis: %v\n", err)
		}
		// does not have reacted
		return false
	}
	return true
}

// SaveReaction saves that an user reacted to an offer
func (s *userService) SaveReaction(uuid string) {
	err := s.rdb.Set(uuid, true, 0).Err()
	if err != nil {
		log.Printf("error when saving reaction to redis: %v\n", err)
	}
}

func buildReactionUUID(user *User, offer corporation.Offer) string {
	return base64.StdEncoding.EncodeToString([]byte(user.Email + offer.Housing.Address))
}
