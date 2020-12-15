package user

import (
	"errors"
	"fmt"
	"log"

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

	MatchOffer(offers corporation.Offer) error
}

// userService represents a PostgreSQL implementation of Service.
type userService struct {
	db                 *gorm.DB
	aesSecret          string
	clientProvider     corporation.ClientProvider
	corporationService corporation.Service
}

func NewService(db *gorm.DB, aesSecret string, clientProvider corporation.ClientProvider, corporationService corporation.Service) Service {
	return &userService{
		db:                 db,
		aesSecret:          aesSecret,
		clientProvider:     clientProvider,
		corporationService: corporationService,
	}
}

func (s *userService) CreateUser(u *User) (*User, error) {
	if u.Email == "" {
		return nil, fmt.Errorf("email is required for creating user")
	}
	_, err := s.GetUser(u.Email)
	if err == nil {
		return nil, fmt.Errorf("error user %s already exists", u.Email)
	}

	if len(u.HousingPreferences.Type) == 0 {
		return nil, fmt.Errorf("housing preferences is required for creating user")
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

	if err := s.db.Model(u).Association("HousingPreferences").Replace(&housingPreferences); err != nil {
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

func (s *userService) MatchOffer(offer corporation.Offer) error {
	// find user with credentials for offer corporation

	// find user with housing preferences matching offer

	// check redis - when reacting store ID, Corporation to check if need to react again

	// react to matching offers

	// store in redis

	// send mail

	// var hasApplied int
	log.Printf("checking %s...\n", offer.Housing.Address)
	// 	if user.MatchPreferences(offer) && user.MatchCriteria(offer) {
	// 		err := client.ReactToOffer(offer)
	// 		if err != nil {
	// 			log.Fatal(err)
	// 			continue
	// 		}

	// 		log.Printf("successfuly applied to %s - view on %s", offer.Housing.Address, offer.URL)
	// 		hasApplied++
	// 	}

	// log.Printf("WoningFinder has applied to %d house(s) on behalf of %s today.", hasApplied, user.Name)

	return nil
}
