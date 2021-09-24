package domijn

import (
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/woningfinder/woningfinder/internal/corporation"
	"github.com/woningfinder/woningfinder/internal/corporation/connector"
)

func (c *client) React(offer corporation.Offer) error {
	reactURL := c.corporation.APIEndpoint.String() + offer.ExternalID

	// parse react error
	var hasReacted error
	c.collector.OnScraped(func(resp *colly.Response) {
		hasReacted = checkReact(string(resp.Body))
	})

	if err := c.collector.PostRaw(reactURL, nil); err != nil {
		return err
	}

	return hasReacted
}

func checkReact(body string) error {
	domijnReactMsg := "We hebben jouw reactie ontvangen."

	if !strings.Contains(string(body), domijnReactMsg) {
		return connector.ErrReactUnknown
	}

	return nil
}
