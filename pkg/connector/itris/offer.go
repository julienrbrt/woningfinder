package itris

import (
	"github.com/gocolly/colly"
	"github.com/woningfinder/woningfinder/internal/corporation"
)

func (c *itrisConnector) GetOffer() ([]corporation.Offer, error) {
	offerURL := c.url.String() + "/woningaanbod/"

	//get offer
	var offers []corporation.Offer
	c.collector.OnHTML("div.aanbodListItems", func(e *colly.HTMLElement) {
		e.ForEach("div.woningaanbod", func(_ int, i *colly.HTMLElement) {
			externalID, _ := i.DOM.Attr("data-aanbod-id")
			address, _ := i.DOM.Attr("data-select-address")

			offer := corporation.Offer{
				ExternalID: externalID,
				Housing: corporation.Housing{
					Address: address,
				},
			}

			offers = append(offers, offer)
		})
	})

	if err := c.collector.Visit(offerURL); err != nil {
		return nil, err
	}

	return offers, nil
}
