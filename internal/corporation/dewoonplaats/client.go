package dewoonplaats

import (
	"github.com/woningfinder/woningfinder/internal/corporation"
	"github.com/woningfinder/woningfinder/pkg/networking"
)

// Client defines De Woonplaats client
type Client interface {
	Login(username, password string) error
	FetchOffer() ([]corporation.Housing, error)
}

type client struct {
	networkingClient networking.Client
}

// NewClient creates a client for De Woonplaats
func NewClient(c networking.Client) Client {
	return &client{
		networkingClient: c,
	}
}
