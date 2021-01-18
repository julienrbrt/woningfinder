package dewoonplaats

import (
	"github.com/woningfinder/woningfinder/internal/corporation"
	"github.com/woningfinder/woningfinder/pkg/logging"
	"github.com/woningfinder/woningfinder/pkg/mapbox"
	"github.com/woningfinder/woningfinder/pkg/networking"
)

type client struct {
	logger           *logging.Logger
	networkingClient networking.Client
	mapboxClient     mapbox.Client
}

// NewClient creates a client for De Woonplaats
func NewClient(logger *logging.Logger, c networking.Client, mapboxClient mapbox.Client) corporation.Client {
	return &client{
		logger:           logger,
		networkingClient: c,
		mapboxClient:     mapboxClient,
	}
}

func (c *client) Send(req networking.Request) (*response, error) {
	// send request to networking client
	resp, err := c.networkingClient.Send(&req)
	if err != nil {
		return nil, err
	}

	var r response
	err = resp.ReadJSONBody(&r)
	if err != nil {
		return nil, err
	}

	return &r, r.Error()
}
