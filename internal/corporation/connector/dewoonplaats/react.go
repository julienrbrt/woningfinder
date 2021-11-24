package dewoonplaats

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienrbrt/woningfinder/internal/corporation"
	"github.com/julienrbrt/woningfinder/pkg/networking"
)

const methodReact = "ReageerOpWoning"

func (c *client) React(offer corporation.Offer) error {
	req, err := reactRequest(offer.ExternalID)
	if err != nil {
		return err
	}

	_, err = c.Send(req)
	return err
}

func reactRequest(id string) (networking.Request, error) {
	req := request{
		ID:     1,
		Method: methodReact,
		Params: []string{
			id,
		},
	}

	body, err := json.Marshal(req)
	if err != nil {
		return networking.Request{}, fmt.Errorf("error while marshaling %v: %w", req, err)
	}

	request := networking.Request{
		Path:   "/woonplaats_digitaal/woonvinder",
		Method: http.MethodPost,
		Body:   bytes.NewBuffer(body),
	}

	return request, nil
}
