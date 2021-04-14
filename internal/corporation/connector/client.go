package connector

import (
	"github.com/woningfinder/woningfinder/internal/corporation"
)

// Client defines a housing corporation client
type Client interface {
	Login(username, password string) error
	GetOffers() ([]corporation.Offer, error)
	React(offer corporation.Offer) error
}
