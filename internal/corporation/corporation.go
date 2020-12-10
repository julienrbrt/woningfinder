package corporation

import (
	"time"

	"gorm.io/gorm"
)

// Corporation defines a housing corporations basic data
// That data is shared between every housing corporations
type Corporation struct {
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       gorm.DeletedAt `gorm:"index"`
	Name            string         `gorm:"primaryKey"`
	URL             string
	Cities          []City            `gorm:"primaryKey;many2many:corporation_cities"`
	SelectionMethod []SelectionMethod `gorm:"many2many:corporation_selection_method"`
}

// City defines a city where a HousingCorporation operates or when an house offer lies
type City struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Name      string         `gorm:"primaryKey"`
	Region    string         `gorm:"primaryKey"`
}

// Offer defines a house or an appartement available in a Housing Corporation
type Offer struct {
	ExternalID           string // identifier of the house at the housing coporation in order to react
	Housing              Housing
	URL                  string
	SelectionMethod      SelectionMethod
	SelectionDate        time.Time
	City                 City
	CanApply, HasApplied bool
}

const (
	House HousingType = iota
	Appartement
	Parking
	Undefined
)

// HousingType defines the type of an Housing (appartement, house)
type HousingType int

// Housing defines an appartement and a house
type Housing struct {
	Type                    HousingType
	Address                 string
	District                string
	EnergieLabel            string
	Price                   float64
	Size                    float64
	Longitude, Latitude     float64
	NumberRoom              int
	NumberBedroom           int
	BuildingYear            int
	HousingAllowance        bool // HousingAllowance defines if the house can get housing allowance
	Garden                  bool
	Garage                  bool
	Elevator                bool
	Balcony                 bool
	AccessibilityWheelchair bool
	AccessibilityScooter    bool
	Attic                   bool
	Historic                bool // Defines if house is an historical monument
	CV                      bool // Defines if the house has a central verwarming
}
