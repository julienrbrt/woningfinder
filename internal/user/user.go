package user

import (
	"gorm.io/gorm"

	"github.com/woningfinder/woningfinder/internal/corporation"
)

// User defines an user of WoningFinder
type User struct {
	gorm.Model
	Name                   string
	Email                  string
	BirthYear              int
	YearlyIncome           int
	FamilySize             int
	HousingPreferences     HousingPreferences       `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	CorporationCredentials []CorporationCredentials `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

// HousingPreferences defines the user preference on a housing
type HousingPreferences struct {
	gorm.Model
	UserID              int                       `gorm:"primaryKey"`
	Type                []corporation.HousingType `gorm:"many2many:housing_preferences_housing_types"`
	MaximumPrice        float64
	City                []corporation.City         `gorm:"many2many:housing_preferences_cities"`
	CityDistrict        []corporation.CityDistrict `gorm:"many2many:housing_preferences_city_districts"`
	NumberBedroom       int
	HasBalcony          bool
	HasGarage           bool
	HasGarden           bool
	HasElevator         bool
	HasHousingAllowance bool
	IsAccessible        bool
}

// CorporationCredentials holds the user credentials to login to an housing corporation
type CorporationCredentials struct {
	gorm.Model
	UserID          int    `gorm:"primaryKey"`
	CorporationName string `gorm:"primaryKey"`
	CorporationURL  string `gorm:"primaryKey"`
	Corporation     corporation.Corporation
	Login           string
	Password        string
}
