package dewoonplaats

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienrbrt/woningfinder/internal/corporation/connector"
	"github.com/julienrbrt/woningfinder/pkg/networking"
)

const methodLogin = "Login"

type loginResponse struct {
	Code    string `json:"code"`
	Success bool   `json:"success"`
}

// Authenticate to De Woonplaats
func (c *client) Login(username, password string) error {
	req, err := loginRequest(username, password)
	if err != nil {
		return err
	}

	resp, err := c.Send(req)
	if err != nil {
		return err
	}

	var response loginResponse
	if err := json.Unmarshal(resp.Result, &response); err != nil {
		return fmt.Errorf("error parsing login result %v: %w", resp.Result, err)
	}

	if !response.Success {
		return fmt.Errorf("error authentication %s: %w", response.Code, connector.ErrAuthFailed)
	}

	return nil
}

func loginRequest(username, password string) (networking.Request, error) {
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
		return networking.Request{}, fmt.Errorf("error while marshaling %v: %w", req, err)
	}

	request := networking.Request{
		Path:   "/wrd/auth",
		Method: http.MethodPost,
		Body:   bytes.NewBuffer(body),
	}

	return request, nil
}
