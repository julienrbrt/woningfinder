package itris

import (
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/woningfinder/woningfinder/internal/corporation"
)

const detailsHousingChildAttr = "li.link a"

func (c *client) GetOffers() ([]corporation.Offer, error) {
	offers := map[string]*corporation.Offer{}

	// add offer
	c.collector.OnHTML("div.aanbodListItems", func(el *colly.HTMLElement) {
		el.ForEach("div.woningaanbod", func(_ int, e *colly.HTMLElement) {
			var offer corporation.Offer
			var err error

			// get offer url
			offer.URL = c.url + e.ChildAttr(detailsHousingChildAttr, "href")
			offer.ExternalID = e.Attr("data-aanbod-id")

			// get housing type
			offer.Housing.Type = c.parseHousingType(e.Text)
			if offer.Housing.Type == corporation.HousingTypeUndefined {
				return
			}

			// get location
			offer.Housing.Address = strings.Title(strings.ToLower(e.ChildAttr(detailsHousingChildAttr, "data-select-address")))
			offer.Housing.City.Name = strings.Title(strings.ToLower(e.Attr("data-plaats")))
			offer.Housing.CityDistrict, err = c.mapboxClient.CityDistrictFromAddress(offer.Housing.Address)
			if err != nil {
				c.logger.Sugar().Warnf("could not get city district of %s: %w", offer.Housing.Address, err)
			}

			offer.SelectionDate, err = time.Parse(layoutTime, e.Attr("data-reactiedatum"))
			if err != nil {
				c.logger.Sugar().Errorf("error while parsing date of %s: %w", offer.Housing.Address, err)
				return
			}

			offer.Housing.Price, err = strconv.ParseFloat(e.Attr("data-prijs"), 16)
			if err != nil {
				c.logger.Sugar().Errorf("error while parsing price of %s: %w", offer.Housing.Address, err)
				return
			}

			offer.Housing.NumberBedroom, err = strconv.Atoi(e.Attr("data-kamers"))
			if err != nil {
				c.logger.Sugar().Errorf("error while parsing number bedroom of %s: %w", offer.Housing.Address, err)
				return
			}

			// create new offer
			offers[offer.URL] = &offer

			// visit offer url
			if err := e.Request.Visit(offer.URL); err != nil {
				c.logger.Sugar().Errorf("error while checking offer details %s: %w", offer.Housing.Address, err)
			}
		})
	})

	// add housing details
	c.collector.OnHTML("div.info-container", func(e *colly.HTMLElement) {
		// find offer
		offerURL := e.Request.URL.String()
		if _, ok := offers[offerURL]; !ok {
			// no offer matching, no details
			return
		}

		c.housingDetailsParser(c.logger, offers[offerURL], e)
	})

	// parse offers
	offerURL := c.url + "/woningaanbod/"
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

func (c *client) parseHousingType(houseType string) corporation.HousingType {
	houseType = strings.ToLower(houseType)

	if strings.Contains(houseType, "appartement") {
		return corporation.HousingTypeAppartement
	}

	if strings.Contains(houseType, "woning") {
		return corporation.HousingTypeHouse
	}

	return corporation.HousingTypeUndefined
}
