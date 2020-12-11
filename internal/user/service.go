package user

import (
	"fmt"

	"gorm.io/gorm/clause"

	"github.com/woningfinder/woningfinder/internal/corporation"
	"gorm.io/gorm"
)

// Service permits to handle the persistence of an user
type Service interface {
	Create(u *User) (*User, error)
	GetUser(email string) (*User, error)
	Delete(u *User) error

	UpdateHousingPreferences(u *User, pref HousingPreferences) error
	GetHousingPreferences(u *User) (HousingPreferences, error)
	DeleteHousingPreferences(u *User) error

	CreateCorporationCredentials(u *User, corporation corporation.Corporation) error
	GetCorporationCredentials(u *User, corporation corporation.Corporation) (CorporationCredentials, error)
	DeleteCorporationCredentials(u *User, corporation corporation.Corporation) error

	// Matcher
	MatchOffer(offers corporation.Offer) error
}

// userService represents a PostgreSQL implementation of Service.
type userService struct {
	db                 *gorm.DB
	corporationService corporation.Service
}

func NewService(db *gorm.DB, corporationService corporation.Service) Service {
	return &userService{
		db:                 db,
		corporationService: corporationService,
	}
}

// Create an user in the database
func (s *userService) Create(u *User) (*User, error) {
	if u.Email == "" {
		return nil, fmt.Errorf("email is required for creating user")
	}

	if len(u.HousingPreferences.Type) == 0 {
		return nil, fmt.Errorf("housing preferences is required for creating user")
	}

	for i, city := range u.HousingPreferences.City {
		newCity, err := s.corporationService.GetCity(city)
		if err != nil {
			return nil, fmt.Errorf("the given city %s is not supported", city.Name)
		}

		u.HousingPreferences.City[i] = *newCity
	}

	result := s.db.FirstOrCreate(u)
	if result.Error != nil {
		return nil, result.Error
	}

	return u, nil
}

// GetUser from an email
func (s *userService) GetUser(email string) (*User, error) {
	var u User
	s.db.Where(&User{Email: email}).First(&u)

	if u.ID == 0 {
		return nil, fmt.Errorf("no user found with the email: %s", email)
	}

	return &u, nil
}

// Delete an user from the database
func (s *userService) Delete(u *User) error {
	result := s.db.Select(clause.Associations).Delete(u)
	return result.Error
}

func (s *userService) UpdateHousingPreferences(u *User, pref HousingPreferences) error {
	return nil
}

func (s *userService) GetHousingPreferences(u *User) (HousingPreferences, error) {
	return HousingPreferences{}, nil
}

func (s *userService) DeleteHousingPreferences(u *User) error {
	return nil
}

func (s *userService) CreateCorporationCredentials(u *User, corporation corporation.Corporation) error {
	return nil
}

func (s *userService) GetCorporationCredentials(u *User, corporation corporation.Corporation) (CorporationCredentials, error) {
	return CorporationCredentials{}, nil
}

func (s *userService) DeleteCorporationCredentials(u *User, corporation corporation.Corporation) error {
	return nil
}

func (s *userService) MatchOffer(offer corporation.Offer) error {
	fmt.Println("--------")
	fmt.Println(offer.Housing.Address)
	fmt.Println("--------")

	// find user with housing preferences matching offer

	// find user with credentials for offer corporation

	// apply to matching offers

	// store in redis

	// send mail

	// var hasApplied int
	// 	log.Printf("checking %s...\n", offer.Housing.Address)
	// 	if user.MatchPreferences(offer) && user.MatchCriteria(offer) {
	// 		err := client.ApplyOffer(offer)
	// 		if err != nil {
	// 			log.Fatal(err)
	// 			continue
	// 		}

	// 		log.Printf("successfuly applied to %s - view on %s", offer.Housing.Address, offer.URL)
	// 		hasApplied++
	// 	}

	// REDIS
	// when applying store ID, Corporation to check if need to apply again
	// log.Printf("WoningFinder has applied to %d house(s) on behalf of %s today.", hasApplied, user.FullName)

	return nil
}
