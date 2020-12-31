package corporation

import (
	"database/sql/driver"
	"time"

	"gorm.io/gorm"
)

const (
	// House is a house type
	House Type = "HOUSE"
	// Appartement is a appartement type
	Appartement = "APPARTEMENT"
	// Undefined is an undefined type (probably parking but unsupported by WoningFinder)
	Undefined = "UNDEFINED"
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

// IsValid confirms that a housing contains all the required information by WoningFinder
func (h *Housing) IsValid() bool {
	if h.Type.Type == Undefined {
		return true
	}

	return h.Type.Type != "" &&
		h.City.Name != "" &&
		h.CityDistrict.CityName == h.City.Name &&
		h.EnergieLabel != "" &&
		h.BuildingYear > 0 &&
		h.Price > 0 &&
		h.NumberBedroom > 0
}
