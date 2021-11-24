package itris

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/julienrbrt/woningfinder/internal/corporation"
	"github.com/julienrbrt/woningfinder/internal/corporation/connector"
)

func (c *client) React(offer corporation.Offer) error {
	// e.g. https://mijn.onshuis.com/mijn-portaal/inschrijvingen/aanmakenWoningreactie/index.xml?reac_csrf_protection=65159aa78108a8370fc71d49a95aa416859b18e1&ogeNr=380000000020135&PublicatieNr=380009660 - it seems that ogeNr is not required
	reactURL := fmt.Sprintf("%s/mijn-portaal/inschrijvingen/aanmakenWoningreactie/index.xml?reac_csrf_protection=%s&PublicatieNr=%s", c.corporation.APIEndpoint.String(), c.itrisCSRFToken, offer.ExternalID)

	// parse react error
	var hasReacted error
	c.collector.OnScraped(func(resp *colly.Response) {
		hasReacted = c.checkReact(string(resp.Body))
	})

	if err := c.collector.Visit(reactURL); err != nil {
		return err
	}

	return hasReacted
}

func (c *client) checkReact(body string) error {
	if !strings.Contains(string(body), "Uw woningreactie is aangemaakt.") {
		return connector.ErrReactUnknown
	}

	return nil
}
