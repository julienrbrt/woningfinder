package corporation

import (
	"database/sql/driver"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/woningfinder/woningfinder/pkg/osm"
	"gorm.io/gorm"
)

// HousingTypeDB permits to create the housing type in the database
var HousingTypeDB = []HousingType{
	{
		Type: House,
	},
	{
		Type: Appartement,
	},
	{
		Type: Parking,
	},
	{
		Type: Undefined,
	},
}

const (
	House       Type = "HOUSE"
	Appartement      = "APPARTEMENT"
	Parking          = "PARKING"
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
	CityID                  int
	City                    City
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

// SetDistrict set the district name from a location
func (h *Housing) SetDistrict() {
	if len(h.City.District) > 0 {
		h.City.District[0].Name = strings.ToLower(h.City.District[0].Name)
		return
	}

	name, err := osm.GetResidential(fmt.Sprintf("%.5f", h.Latitude), fmt.Sprintf("%.5f", h.Longitude))
	if err != nil {
		log.Printf(fmt.Errorf("error getting district from %s: %w", h.Address, err).Error())
	}

	if name == "" {
		h.City.District = nil
		return
	}

	h.City.District = append(h.City.District, District{Name: name})
}
