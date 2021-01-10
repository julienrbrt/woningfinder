package dewoonplaats

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/woningfinder/woningfinder/internal/corporation"

	"github.com/woningfinder/woningfinder/pkg/networking"
)

const methodLogin = "Login"

type loginResult struct {
	Code     string `json:"code"`
	Success  bool   `json:"success"`
	Userinfo struct {
		Name     string `json:"fullname"`
		Gender   string `json:"geslacht"`
		Age      int    `json:"leeftijd"`
		Postcode string `json:"postcode"`
	} `json:"userinfo"`
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

	var result loginResult
	if err := json.Unmarshal(resp.Result, &result); err != nil {
		return fmt.Errorf("error parsing login result %v: %w", string(resp.Result), err)
	}

	if !result.Success {
		return fmt.Errorf("error authentication %s: %w", result.Code, corporation.ErrAuthFailed)
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
