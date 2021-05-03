package zig

import (
	"encoding/json"
	"fmt"

	"github.com/woningfinder/woningfinder/internal/corporation/connector"
	"github.com/woningfinder/woningfinder/pkg/logging"
	"github.com/woningfinder/woningfinder/pkg/mapbox"
	"github.com/woningfinder/woningfinder/pkg/networking"
)

type client struct {
	logger           *logging.Logger
	networkingClient networking.Client
	mapboxClient     mapbox.Client
}

// NewClient creates a client for Zig
func NewClient(logger *logging.Logger, c networking.Client, mapboxClient mapbox.Client) connector.Client {
	return &client{
		logger:           logger,
		networkingClient: c,
		mapboxClient:     mapboxClient,
	}
}

func (c *client) Send(req networking.Request) (json.RawMessage, error) {
	// add header
	if req.Headers == nil {
		req.Headers = map[string]string{}
	}
	req.Headers["Content-Type"] = "application/x-www-form-urlencoded"

	// send request to networking client
	resp, err := c.networkingClient.Send(&req)
	if err != nil {
		return nil, err
	}

	var rawResponse json.RawMessage
	if err := resp.ReadJSONBody(&rawResponse); err != nil {
		// fallback on response error handling
		type responseError struct {
			Err string `json:"sMessage"`
		}

		var response responseError
		if err := resp.ReadJSONBody(&response); err != nil {
			return nil, fmt.Errorf("error unmarshaling zig response %v: %w", resp, err)
		}

		return nil, fmt.Errorf("zig error reponse: %v", response.Err)
	}

	return rawResponse, nil
}
