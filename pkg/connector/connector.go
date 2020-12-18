package connector

import "github.com/woningfinder/woningfinder/internal/corporation"

// Connector specifies information about the itris running instance
type Connector interface {
	GetOffer() ([]corporation.Offer, error)
}
