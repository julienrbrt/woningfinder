package dewoonplaats

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/woningfinder/woningfinder/internal/corporation"
	"github.com/woningfinder/woningfinder/pkg/networking"
)

const methodApply = "ReageerOpWoning"

func (c *client) ApplyOffer(offer corporation.Offer) error {
	req, err := applyRequest(offer.ExternalID)
	if err != nil {
		return err
	}

	_, err = c.Send(req)
	return err
}

func applyRequest(id string) (networking.Request, error) {
	req := request{
		ID:     1,
		Method: methodApply,
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
