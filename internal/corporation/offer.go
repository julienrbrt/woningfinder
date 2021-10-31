package corporation

import "net/url"

// Offers defines a list of offer belonging to a housing corporation
type Offers struct {
	Corporation Corporation
	Offer       []Offer
}

// Offer defines a house or an appartement available in a housing corporation
type Offer struct {
	Housing       Housing
	URL           string
	RawPictureURL *url.URL
	// identifier of the house at the housing coporation in order to react
	ExternalID                   string
	MinFamilySize, MaxFamilySize int
	MinAge, MaxAge               int
	MinimumIncome, MaximumIncome int
}
