package corporation

import (
	"net/url"
	"time"
)

const (
	// SelectionRandom selects a candidate from an offer randomly
	SelectionRandom SelectionMethod = iota
	// SelectionFirstComeFirstServed selects first candidate that reacted to an offer
	SelectionFirstComeFirstServed
	// SelectionRegistrationDate selects the candidate that registered the first in the housing corporation in the offer drawing
	SelectionRegistrationDate
)

// SelectionMethod defines the selection method used for a Housing Corporation in an Offer
// There is 3 supported SelectionMethod: SelectionRandom, SelectionFirstComeFirstServed, SelectionRegistrationDate
type SelectionMethod int

// Corporation defines a housing corporations basic data
// That data is shared between every housing corporations
type Corporation struct {
	Name            string
	URL             url.URL
	Cities          []City
	SelectionMethod []SelectionMethod
}

// City defines a city where a HousingCorporation operates or when an house offer lies
type City struct {
	Name string
}

// Location defines the Longitude and Latidude a city, or offer
type Location struct {
	Longitude, Latitude float64
}

// Offer defines a house or an appartement available in a Housing Corporation
type Offer struct {
	ExternalID           string // identifier of the house at the housing coporation in order to react
	Housing              Housing
	URL                  *url.URL
	PictureURL           *url.URL
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
	Location                Location
	Address                 string
	District                string
	EnergieLabel            string
	Price                   float64
	Size                    float64
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
