package onshuis

import (
	"net/url"

	"github.com/woningfinder/woningfinder/internal/corporation"
	"github.com/woningfinder/woningfinder/pkg/networking"
)

// Host defines the OnsHuis offers domain
var Host = &url.URL{Scheme: "https", Host: "mijn.onshuis.com", Path: "/apps/com.itris.klantportaal/"}

type client struct {
	networkingClient networking.Client
}

// NewClient creates a client for OnsHuis
func NewClient(c networking.Client) corporation.Client {
	return &client{
		networkingClient: c,
	}
}

func (c *client) Login(username, password string) error {
	return nil
}

func (c *client) FetchOffer(minimumPrice float64) ([]corporation.Offer, error) {
	return nil, nil
}

func (c *client) ApplyOffer(offer corporation.Offer) error {
	return nil
}