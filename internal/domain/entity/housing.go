package entity

import (
	"database/sql/driver"
)

const (
	// HousingTypeHouse is a house type
	HousingTypeHouse Type = "HOUSE"
	// HousingTypeAppartement is a appartement type
	HousingTypeAppartement = "APPARTEMENT"
	// HousingTypeUndefined is an undefined type (probably parking but unsupported by WoningFinder)
	HousingTypeUndefined = "UNDEFINED"
)

// Type defines the type of an HousingType (appartement, house)
type Type string

// Scan implements the Scanner interface from reading from the database
func (u *Type) Scan(value interface{}) error {
	*u = Type(string(value.([]byte)))
	return nil
}

// Value implements the Valuer interface for the storing in the database
func (u Type) Value() (driver.Value, error) {
	return string(u), nil
}

// HousingType is the database representation of (Housing)Type
type HousingType struct {
	Type Type `pg:",pk"`
}

// Exists verify that the given housing type is supported
// Note HousingTypeUndefined is not returned as exists as it cannot be selected by a user
func (h *HousingType) Exists() bool {
	return h.Type == HousingTypeHouse || h.Type == HousingTypeAppartement
}

// Housing defines an appartement and a house
type Housing struct {
	Type                HousingType
	Address             string
	City                City
	CityDistrict        string
	EnergieLabel        string
	Price               float64
	Size                float64
	Longitude, Latitude float64
	NumberRoom          int
	NumberBedroom       int
	BuildingYear        int
	HousingAllowance    bool // HousingAllowance defines if the house can get housing allowance
	Garden              bool
	Garage              bool
	Elevator            bool
	Balcony             bool
	Attic               bool
	Accessible          bool // Assessible defines if the house is accessible for handicapt people
}
