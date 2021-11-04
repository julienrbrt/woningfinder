package woningnet

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/woningfinder/woningfinder/internal/corporation"
	"github.com/woningfinder/woningfinder/internal/corporation/connector"
)

func (c *client) React(offer corporation.Offer) error {
	reactURL := c.corporation.URL + "/Zoeken/Eenheiddetails/Reageren/" + offer.ExternalID
	reactRequest := map[string][]byte{
		"__RequestVerificationToken": nil,
		"formValidation":             []byte("on"),
	}

	// we clone the collector in order to send the request in a different collector than we visit the url
	// this permit to do not fall in an infinite loop
	collector := c.collector.Clone()

	// behave like a web browser
	collector.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Sec-Fetch-Site", "same-origin")
		r.Headers.Set("Sec-Fetch-Mode", "navigate")
		r.Headers.Set("Sec-Fetch-User", "?1")
		r.Headers.Set("Sec-Fetch-Dest", "document")
		r.Headers.Set("Upgrade-Insecure-Requests", "1")
	})

	// parse react page
	c.collector.OnHTML("form", func(el *colly.HTMLElement) {
		// fill react form
		reactRequest["__RequestVerificationToken"] = []byte(el.ChildAttr("input[name=__RequestVerificationToken]", "value"))

		// send request
		collector.PostMultipart(reactURL, reactRequest)
	})

	// parse react error (from second collector)
	var hasReacted error
	collector.OnScraped(func(resp *colly.Response) {
		hasReacted = c.checkReact(string(resp.Body))
	})

	// visit login page
	if err := c.collector.Visit(reactURL); err != nil {
		return fmt.Errorf("error visiting react page: %w", err)
	}

	return hasReacted
}

func (c *client) checkReact(body string) error {
	if !strings.Contains(body, "Uw reactie is opgeslagen") {
		return connector.ErrReactUnknown
	}

	return nil
}
