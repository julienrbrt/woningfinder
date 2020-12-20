package corporation

import (
	"database/sql/driver"
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

const (
	// SelectionRandom selects a candidate from an offer randomly
	SelectionRandom Method = "SELECTION_RANDOM"
	// SelectionFirstComeFirstServed selects first candidate that reacted to an offer
	SelectionFirstComeFirstServed = "SELECTION_FIRST_COME_FIRST_SERVED"
	// SelectionRegistrationDate selects the candidate that registered the first in the housing corporation in the offer drawing
	SelectionRegistrationDate = "SELECTION_REGISTRATION_DATE"
)

// Method defines the selection method used for a Housing Corporation in an Offer
// There is 3 supported Method: SelectionRandom, SelectionFirstComeFirstServed, SelectionRegistrationDate
type Method string

// Scan implements the Scanner interface from reading from the database
func (u *Method) Scan(value interface{}) error {
	*u = Method(value.(string))
	return nil
}

// Value implements the Valuer interface for the storing in the database
func (u Method) Value() (driver.Value, error) {
	return string(u), nil
}

// SelectionMethod is the database representation of Method
type SelectionMethod struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Method    Method         `gorm:"primaryKey"`
}

// Offer defines a house or an appartement available in a Housing Corporation
type Offer struct {
	Housing                      Housing
	SelectionMethod              SelectionMethod
	SelectionDate                time.Time
	URL                          string
	ExternalID                   string // identifier of the house at the housing coporation in order to react
	MinIncome, MaxIncome         int
	MinFamilySize, MaxFamilySize int
	MinAge, MaxAge               int
}

// OfferList defines a list of offer belonging to one corporation
type OfferList struct {
	Corporation Corporation
	Offer       []Offer
}
