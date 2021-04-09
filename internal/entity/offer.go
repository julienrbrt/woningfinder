package entity

import (
	"time"
)

// Offer defines a house or an appartement available in a Housing Corporation
type Offer struct {
	Housing                      Housing
	SelectionMethod              SelectionMethod
	SelectionDate                time.Time
	URL                          string
	ExternalID                   string // identifier of the house at the housing coporation in order to react
	MinFamilySize, MaxFamilySize int
	MinAge, MaxAge               int
}

// OfferList defines a list of offer belonging to one corporation
type OfferList struct {
	Corporation Corporation
	Offer       []Offer
}
