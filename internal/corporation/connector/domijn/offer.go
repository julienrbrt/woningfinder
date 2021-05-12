package domijn

import (
	"fmt"

	"github.com/gocolly/colly/v2"
	"github.com/woningfinder/woningfinder/internal/corporation"
)

func (c *client) GetOffers() ([]corporation.Offer, error) {
	offers := map[string]*corporation.Offer{}

	// add offer
	c.collector.OnHTML("#housingWrapper", func(el *colly.HTMLElement) {
		el.ForEach("article", func(_ int, e *colly.HTMLElement) {
			var offer corporation.Offer

			// get offer url
			offer.URL = c.url + e.ChildAttr("div.image-wrapper a", "href")

			fmt.Println(e)
			fmt.Println(offer)

			// create new offer
			offers[offer.URL] = &offer

			// visit offer url
			if err := e.Request.Visit(offer.URL); err != nil {
				c.logger.Sugar().Errorf("itrs connector: error while checking offer details %s: %w", offer.Housing.Address, err)
			}
		})
	})

	// parse offers
	offerURL := c.url + "/woningaanbod/?advertisementType=rent"
	if err := c.collector.Visit(offerURL); err != nil {
		return nil, err
	}

	// get all offers as array
	var offerList []corporation.Offer
	for _, offer := range offers {
		offerList = append(offerList, *offer)
	}

	return offerList, nil
}
