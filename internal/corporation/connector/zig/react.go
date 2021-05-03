package zig

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/woningfinder/woningfinder/internal/corporation"
	"github.com/woningfinder/woningfinder/pkg/networking"
)

func (c *client) React(offer corporation.Offer) error {
	hash, err := c.getReactConfiguration()
	if err != nil {
		return err
	}

	externalID := strings.Split(offer.ExternalID, externalIDSeperator)
	resp, err := c.Send(reactRequest(externalID[0], externalID[1], hash))

	var result struct {
		Success bool `json:"success"`
	}

	if err := json.Unmarshal(resp, &result); err != nil {
		return fmt.Errorf("error parsing react result %v: %w", string(resp), err)
	}

	if !result.Success {
		return fmt.Errorf("reaction to %s not successful", offer.Housing.Address)
	}

	return err
}

func reactRequest(assignmentID, dwellingID, hash string) networking.Request {
	body := url.Values{}
	body.Add("__id__", "Portal_Form_SubmitOnly")
	body.Add("__hash__", hash)
	body.Add("add", assignmentID)
	body.Add("dwellingID", dwellingID)

	request := networking.Request{
		Path:   "/portal/object/frontend/react/format/json",
		Method: http.MethodPost,
		Body:   strings.NewReader(body.Encode()),
	}

	return request
}

func (c *client) getReactConfiguration() (string, error) {
	resp, err := c.Send(networking.Request{
		Path:   "/portal/object/frontend/getreageerconfiguration/format/json",
		Method: http.MethodGet,
	})
	if err != nil {
		return "", err
	}

	var result struct {
		Reactform struct {
			Elements struct {
				Hash struct {
					Initialdata string `json:"initialData"`
				} `json:"__hash__"`
			} `json:"elements"`
		} `json:"reageerConfiguration"`
	}

	if err := json.Unmarshal(resp, &result); err != nil {
		return "", fmt.Errorf("error parsing react configuration result %v: %w", string(resp), err)
	}

	return result.Reactform.Elements.Hash.Initialdata, nil
}
