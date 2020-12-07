package dewoonplaats

import (
	"fmt"

	"github.com/woningfinder/woningfinder/internal/corporation"
	"github.com/woningfinder/woningfinder/pkg/networking"
)

// Client defines De Woonplaats client
type Client interface {
	Login(username, password string) error
	FetchOffer(minimumPrice int) ([]corporation.Offer, error)
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
