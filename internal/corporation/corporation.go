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
	Name     string
	Location Location
}

// Location defines the Longitude and Latidude a city, or offer
type Location struct {
	Longitude, Latitude float64
}

// Offer defines a house or an appartement available in a Housing Corporation
type Offer struct {
	Housing            Housing
	URL                url.URL
	StartDate, EndDate time.Time
	SelectionMethod    SelectionMethod
	City               City
}
