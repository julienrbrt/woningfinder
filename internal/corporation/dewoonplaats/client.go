package dewoonplaats

import (
	"encoding/json"
	"fmt"

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

func (c *client) Send(req networking.Request) (response, error) {
	resp, err := c.networkingClient.Send(&req)
	if err != nil {
		return response{}, fmt.Errorf("request %v has given an error: %w", req, err)
	}

	respBody, err := resp.CopyBody()
	if err != nil {
		return response{}, fmt.Errorf("error while copying body response %v: %w", resp, err)
	}

	var r response
	if err := json.Unmarshal(respBody, &r); err != nil {
		return response{}, fmt.Errorf("failed unmarshaling response %v: %w", resp, err)
	}

	// check for response error
	if r.Error() != nil {
		return response{}, r.Error()
	}

	return r, nil
}
