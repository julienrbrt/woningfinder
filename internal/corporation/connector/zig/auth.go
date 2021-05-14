package zig

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/woningfinder/woningfinder/internal/corporation/connector"
	"github.com/woningfinder/woningfinder/pkg/networking"
)

type loginResponse struct {
	Messages []interface{} `json:"messages"`
	Success  bool          `json:"success"`
	Formhash string        `json:"formHash"`
}

// Authenticate to Zig
func (c *client) Login(username, password string) error {
	hash, err := c.getLoginConfiguration()
	if err != nil {
		return err
	}

	resp, err := c.Send(loginRequest(username, password, hash))
	if err != nil {
		return err
	}

	var response loginResponse
	if err := json.Unmarshal(resp, &response); err != nil {
		return fmt.Errorf("error parsing login result %v: %w", string(resp), err)
	}

	if !response.Success || len(response.Messages) > 0 {
		return fmt.Errorf("error authentication %s: %w", response.Messages, connector.ErrAuthFailed)
	}

	return nil
}

func (c *client) getLoginConfiguration() (string, error) {
	resp, err := c.Send(networking.Request{
		Path:   "/portal/account/frontend/getloginconfiguration/format/json",
		Method: http.MethodGet,
	})
	if err != nil {
		return "", err
	}

	var result struct {
		Loginform struct {
			Elements struct {
				Hash struct {
					Initialdata string `json:"initialData"`
				} `json:"__hash__"`
			} `json:"elements"`
		} `json:"loginForm"`
	}

	if err := json.Unmarshal(resp, &result); err != nil {
		return "", fmt.Errorf("error parsing login configuration result %v: %w", string(resp), err)
	}

	return result.Loginform.Elements.Hash.Initialdata, nil
}

func loginRequest(username, password, hash string) networking.Request {
	body := url.Values{}
	body.Add("__hash__", hash)
	body.Add("__id__", "Account_Form_LoginFrontend")
	body.Add("username", username)
	body.Add("password", password)

	request := networking.Request{
		Path:   "/portal/account/frontend/loginbyservice/format/json",
		Method: http.MethodPost,
		Body:   strings.NewReader(body.Encode()),
	}

	return request
}
