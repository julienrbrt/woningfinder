package corporation

import (
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Service permits to handle the persistence of a corporation
type Service interface {
	Create(corporation *[]Corporation) (*[]Corporation, error)
	CreateHousingType(housingTypes *[]HousingType) (*[]HousingType, error)

	GetCorporations() (*[]Corporation, error)
	GetCity(city City) (*City, error)
}

// corporationService represents a PostgreSQL implementation of Service.
type corporationService struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) Service {
	return &corporationService{
		db: db,
	}
}

func (s *corporationService) Create(corporations *[]Corporation) (*[]Corporation, error) {
	// creates the corporation - on data changes update it
	result := s.db.Clauses(clause.OnConflict{UpdateAll: true}).Create(&corporations)
	if result.Error != nil {
		return nil, result.Error
	}

	return corporations, nil
}

func (s *corporationService) CreateHousingType(housingTypes *[]HousingType) (*[]HousingType, error) {
	// creates housing types
	result := s.db.Clauses(clause.OnConflict{UpdateAll: true}).Create(&housingTypes)
	if result.Error != nil {
		return nil, result.Error
	}

	return housingTypes, nil
}

func (s *corporationService) GetCorporations() (*[]Corporation, error) {
	var corporations []Corporation
	result := s.db.Find(&corporations)
	if result.Error != nil {
		return nil, result.Error
	}

	return &corporations, nil
}

func (s *corporationService) GetCity(city City) (*City, error) {
	var c City
	s.db.Where(city).First(&c)

	if c.Name == "" {
		return nil, fmt.Errorf("no city found with the name: %s", city.Name)
	}

	return &c, nil
}
