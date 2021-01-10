package entity

import (
	"net/url"
	"time"

	"gorm.io/gorm"
)

// Corporation defines a housing corporations basic data
// That data is shared between every housing corporations
type Corporation struct {
	APIEndpoint     *url.URL `gorm:"-"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       gorm.DeletedAt    `gorm:"index"`
	Name            string            `gorm:"primaryKey"`
	URL             string            `gorm:"primaryKey"`
	Cities          []City            `gorm:"many2many:corporations_cities"`
	SelectionMethod []SelectionMethod `gorm:"many2many:corporations_selection_method"`
	SelectionTime   time.Time
}

// City defines a city where a housing corporation operates or when an house offer lies
type City struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Name      string         `gorm:"primaryKey"`
	District  []CityDistrict `gorm:"foreignKey:CityName"`
}

// CityDistrict the district of a city
type CityDistrict struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	CityName  string         `gorm:"primaryKey"`
	Name      string         `gorm:"primaryKey"`
}
