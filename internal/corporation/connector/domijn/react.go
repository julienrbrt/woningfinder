package domijn

import (
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/julienrbrt/woningfinder/internal/corporation"
	"github.com/julienrbrt/woningfinder/internal/corporation/connector"
)

func (c *client) React(offer corporation.Offer) error {
	reactURL := c.corporation.APIEndpoint.String() + offer.ExternalID

	// parse react error
	var hasReacted error
	c.collector.OnScraped(func(resp *colly.Response) {
		hasReacted = c.checkReact(string(resp.Body))
	})

	if err := c.collector.PostRaw(reactURL, nil); err != nil {
		return err
	}

	return hasReacted
}

func (c *client) checkReact(body string) error {
	if !strings.Contains(string(body), "We hebben jouw reactie ontvangen.") {
		return connector.ErrReactUnknown
	}

	return nil
}
