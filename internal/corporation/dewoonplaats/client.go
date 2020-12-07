package dewoonplaats

import (
	"fmt"
	"net/url"

	"github.com/woningfinder/woningfinder/internal/corporation"
	"github.com/woningfinder/woningfinder/pkg/networking"
)

// Host defines the De Woonplaats API domain
var Host = &url.URL{Scheme: "https", Host: "www.dewoonplaats.nl", Path: "/wh_services"}

type client struct {
	networkingClient networking.Client
}

// NewClient creates a client for De Woonplaats
func NewClient(c networking.Client) corporation.Client {
	return &client{
		networkingClient: c,
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