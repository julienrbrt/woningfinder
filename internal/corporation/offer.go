package corporation

import (
	"time"
)

// Offers defines a list of offer belonging to a housing corporation
type Offers struct {
	Corporation Corporation
	Offer       []Offer
}

// Offer defines a house or an appartement available in a housing corporation
type Offer struct {
	Housing                      Housing
	SelectionMethod              SelectionMethod
	SelectionDate                time.Time
	URL                          string
	ExternalID                   string // identifier of the house at the housing coporation in order to react
	MinFamilySize, MaxFamilySize int
	MinAge, MaxAge               int
}