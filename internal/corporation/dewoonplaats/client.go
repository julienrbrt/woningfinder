package dewoonplaats

import (
	"fmt"

	"go.uber.org/zap"

	"github.com/woningfinder/woningfinder/internal/corporation"
	"github.com/woningfinder/woningfinder/pkg/mapbox"
	"github.com/woningfinder/woningfinder/pkg/networking"
)

type client struct {
	logger           *zap.Logger
	networkingClient networking.Client
	mapboxClient     mapbox.Client
}

// NewClient creates a client for De Woonplaats
func NewClient(logger *zap.Logger, c networking.Client, mapboxClient mapbox.Client) corporation.Client {
	return &client{
		logger:           logger,
		networkingClient: c,
		mapboxClient:     mapboxClient,
	}
}

func (c *client) Send(req networking.Request) (response, error) {
	// send request to networking client
	resp, err := c.networkingClient.Send(&req)
	if err != nil {
		return response{}, fmt.Errorf("request %v has given an error: %w", req, err)
	}

	var r response
	err = resp.ReadJSONBody(&r)
	if err != nil {
		return response{}, err
	}

	return r, r.Error()
}
