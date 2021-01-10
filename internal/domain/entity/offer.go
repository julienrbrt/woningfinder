package entity

import (
	"time"

	"gorm.io/gorm"
)

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

// OfferMatch defines an offer that matched with an user
// It is used to determined to which offer WoningFinder has applied
type OfferMatch struct {
	gorm.Model
	UserID          int `gorm:"primaryKey"`
	HousingAddress  string
	CorporationName string
	OfferURL        string
}
