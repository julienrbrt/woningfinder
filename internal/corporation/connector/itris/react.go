package itris

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/woningfinder/woningfinder/internal/corporation"
	"github.com/woningfinder/woningfinder/internal/corporation/connector"
)

func (c *client) React(offer corporation.Offer) error {
	// e.g. https://mijn.onshuis.com/mijn-portaal/inschrijvingen/aanmakenWoningreactie/index.xml?reac_csrf_protection=65159aa78108a8370fc71d49a95aa416859b18e1&ogeNr=380000000020135&PublicatieNr=380009660 - it seems that ogeNr is not required
	reactURL := fmt.Sprintf("%s/mijn-portaal/inschrijvingen/aanmakenWoningreactie/index.xml?reac_csrf_protection=%s&PublicatieNr=%s", c.url, c.itrisCSRFToken, offer.ExternalID)

	// parse react error
	var hasReacted error
	c.collector.OnScraped(func(resp *colly.Response) {
		hasReacted = checkReact(string(resp.Body))
	})

	if err := c.collector.Visit(reactURL); err != nil {
		return err
	}

	return hasReacted
}

func checkReact(body string) error {
	itrisReactMsg := "Uw woningreactie is aangemaakt."

	if !strings.Contains(string(body), itrisReactMsg) {
		return connector.ErrReactUnknown
	}

	return nil
}
