package corporation

import (
	"database/sql/driver"
	"fmt"
	"time"

	"github.com/woningfinder/woningfinder/pkg/osm"
	"gorm.io/gorm"
)

const (
	House       Type = "HOUSE"
	Appartement      = "APPARTEMENT"
	Undefined        = "UNDEFINED"
)

// Type defines the type of an HousingType (appartement, house)
type Type string

// Scan implements the Scanner interface from reading from the database
func (u *Type) Scan(value interface{}) error {
	*u = Type(value.(string))
	return nil
}

// Value implements the Valuer interface for the storing in the database
func (u Type) Value() (driver.Value, error) {
	return string(u), nil
}

// HousingType is the database representation of (Housing)Type
type HousingType struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Type      Type           `gorm:"primaryKey"`
}

// Housing defines an appartement and a house
type Housing struct {
	Type                    HousingType
	Address                 string
	City                    City
	CityDistrict            CityDistrict
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
}

// SetCityDistrict set the district name from a location
func (h *Housing) SetCityDistrict() error {
	if h.CityDistrict.Name != "" {
		return nil
	}

	name, err := osm.GetResidential(fmt.Sprintf("%.5f", h.Latitude), fmt.Sprintf("%.5f", h.Longitude))
	if err != nil {
		return fmt.Errorf("error getting district from %s: %w", h.Address, err)
	}

	h.CityDistrict.Name = name
	return nil
}
