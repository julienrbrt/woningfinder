package ikwilhuren

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/julienrbrt/woningfinder/internal/corporation/connector"
	"github.com/julienrbrt/woningfinder/pkg/networking"
)

type loginResponse struct {
	Username            string `json:"username"`
	IsPayedAccount      bool   `json:"is_payed_account"`
	StartSessionAt      string `json:"start_session_at"`
	ExpirationSessionAt string `json:"expiration_session_at"`
}

func (c *client) Login(username, password string) error {
	resp, err := c.Send(loginRequest(username, password))
	if err != nil {
		// login failed if returns 403
		if netErr, ok := networking.AsNetworkingError(err); ok {
			if netErr.Response().StatusCode == http.StatusForbidden {
				return fmt.Errorf("error authentication: %w", connector.ErrAuthFailed)
			}
		}

		return err
	}

	var response loginResponse
	if err := json.Unmarshal(resp, &response); err != nil {
		return fmt.Errorf("error parsing login result %v: %w", resp, err)
	}

	return nil
}

func loginRequest(username, password string) networking.Request {
	body := url.Values{}
	body.Add("isAjax", "true")
	body.Add("username", username)
	body.Add("password", password)

	request := networking.Request{
		Path:   "/user/login",
		Method: http.MethodPost,
		Body:   strings.NewReader(body.Encode()),
	}

	return request
}
