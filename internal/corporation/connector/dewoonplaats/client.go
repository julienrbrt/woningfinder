package dewoonplaats

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

// NewClient creates a client for De Woonplaats
func NewClient(logger *logging.Logger, c networking.Client, mapboxClient mapbox.Client) connector.Client {
	return &client{
		logger:           logger,
		networkingClient: c,
		mapboxClient:     mapboxClient,
	}
}

// request builds a De Woonplaats request
type request struct {
	ID     int         `json:"id"`
	Method string      `json:"method"`
	Params interface{} `json:"params"`
}

// response corresponds to a De Woonplaats response
type response struct {
	Err    interface{}     `json:"error"`
	ID     int             `json:"id"`
	Result json.RawMessage `json:"result"`
}

func (r *response) Error() error {
	if r.Err != nil {
		return fmt.Errorf("de woonplaats error reponse: %v", r.Err.(string))
	}
	return nil
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
