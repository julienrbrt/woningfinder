package ikwilhuren

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/julienrbrt/woningfinder/internal/corporation"
	"github.com/julienrbrt/woningfinder/internal/corporation/connector"
	"github.com/julienrbrt/woningfinder/pkg/networking"
)

type reactResponse struct {
	Success       bool   `json:"success"`
	RentalDealsID string `json:"rentalDealsId"`
}

func (c *client) React(offer corporation.Offer) error {
	resp, err := c.Send(reactRequest(offer.ExternalID))
	if err != nil {
		return err
	}

	var response reactResponse
	if err := json.Unmarshal(resp, &response); err != nil {
		return fmt.Errorf("error parsing react result %v: %w", resp, err)
	}

	if !response.Success {
		return connector.ErrReactUnknown
	}

	return err
}

func reactRequest(id string) networking.Request {
	body := url.Values{}
	body.Add("object", id)
	body.Add("day", "")
	body.Add("time", "")

	request := networking.Request{
		Path:   "/api/make-rental-deal/method/registered-lead/",
		Method: http.MethodPost,
		Body:   strings.NewReader(body.Encode()),
	}

	return request
}
