package dewoonplaats

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/woningfinder/woningfinder/pkg/networking"
)

// Authenticate to De Woonplaats
func (c *client) Login(username, password string) error {
	req := request{
		ID:     1,
		Method: methodLogin,
		Params: []interface{}{
			"https://www.dewoonplaats.nl/mijn-woonplaats/",
			username,
			password,
			false,
		},
	}

	body, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("error while marshaling %v: %w", req, err)
	}

	request := networking.Request{
		Path:   "wh_services/wrd/auth",
		Method: http.MethodPost,
		Body:   bytes.NewBuffer(body),
	}

	resp, err := c.networkingClient.Send(&request)
	if err != nil {
		return fmt.Errorf("request %v has given an error: %w", req, err)
	}

	// TODO
	b, _ := resp.CopyBody()
	fmt.Println(string(b))

	return nil
}
