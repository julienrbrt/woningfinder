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
				c.logger.Sugar().Warnf("itrs connector: could not get city district of %s: %w", offer.Housing.Address, err)
			}

			offer.SelectionDate, err = time.Parse(layoutTime, e.Attr("data-reactiedatum"))
			if err != nil {
				c.logger.Sugar().Errorf("itrs connector: error while parsing date of %s: %w", offer.Housing.Address, err)
				return
			}

			offer.Housing.Price, err = strconv.ParseFloat(e.Attr("data-prijs"), 16)
			if err != nil {
				c.logger.Sugar().Errorf("itrs connector: error while parsing price of %s: %w", offer.Housing.Address, err)
				return
			}

			offer.Housing.NumberBedroom, err = strconv.Atoi(e.Attr("data-kamers"))
			if err != nil {
				c.logger.Sugar().Errorf("itrs connector: error while parsing number bedroom of %s: %w", offer.Housing.Address, err)
				return
			}

			// create new offer
			offers[offer.URL] = &offer

			// visit offer url
			if err := e.Request.Visit(offer.URL); err != nil {
				c.logger.Sugar().Errorf("itrs connector: error while checking offer details %s: %w", offer.Housing.Address, err)
			}
		})
	})

	// add housing details (1/2)
	c.collector.OnHTML("div.info-container", func(e *colly.HTMLElement) {
		// find offer
		offerURL := e.Request.URL.String()
		if _, ok := offers[offerURL]; !ok {
			// no offer matching, no details
			return
		}

		c.getHousingInfo(offers[offerURL], e)
	})

	// add housing details (2/2)
	c.collector.OnHTML("#icons-container ul", func(e *colly.HTMLElement) {
		// find offer
		offerURL := e.Request.URL.String()
		if _, ok := offers[offerURL]; !ok {
			// no offer matching, no details
			return
		}

		c.getHousingDetails(offers[offerURL], e)
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

func (c *client) getHousingInfo(offer *corporation.Offer, e *colly.HTMLElement) {
	// add selection method
	offer.SelectionMethod = corporation.SelectionRandom // TODO to fix - all houses from onshuis are random

	// add housing size
	e.ForEach("#oppervlaktes-page div.infor-wrapper", func(_ int, el *colly.HTMLElement) {
		// increase size
		roomSize, err := strconv.ParseFloat(strings.ReplaceAll(strings.Trim(el.Text, " m2"), ",", "."), 64)
		if err != nil {
			return
		}
		offer.Housing.Size += roomSize
	})

	// add energie label
	energieLabel := e.ChildText("#Woning-page strong.tag-text")
	if energieLabel != "" {
		offer.Housing.EnergyLabel = energieLabel
	}

	// add building year
	e.ForEach("div.infor-wrapper", func(_ int, el *colly.HTMLElement) {
		buildingYear, err := strconv.Atoi(el.Text)
		if err != nil {
			return
		}
		if buildingYear > 1800 { // random building year so high that it cannot be a number of room
			offer.Housing.BuildingYear = buildingYear
		}
	})

	// part of housing details can be found in housing description (attic and accessibility)
	dom, err := e.DOM.Html()
	if err != nil {
		c.logger.Sugar().Warnf("unable to get details for %s on %s", offer.Housing.Address, offer.URL)
		return
	}
	// add attic
	offer.Housing.Attic = strings.Contains(dom, "zolder")
	// add accessible
	offer.Housing.Accessible = strings.Contains(dom, "toegankelijk")
}

func (c *client) getHousingDetails(offer *corporation.Offer, e *colly.HTMLElement) {

	// add housing details
	e.ForEach("li", func(_ int, el *colly.HTMLElement) {
		switch el.Text {
		case "Balkon":
			offer.Housing.Balcony = el.DOM.HasClass("yes")
		case "Berging":
			offer.Housing.Attic = el.DOM.HasClass("yes")
		case "Garage":
			offer.Housing.Garage = el.DOM.HasClass("yes")
		case "Lift":
			offer.Housing.Elevator = el.DOM.HasClass("yes")
		case "Tuin":
			offer.Housing.Garden = el.DOM.HasClass("yes")
		}
	})
}

func (c *client) parseHousingType(houseType string) corporation.HousingType {
	houseType = strings.ToLower(houseType)

	if strings.Contains(houseType, "appartement") || strings.Contains(houseType, "penthouse") {
		return corporation.HousingTypeAppartement
	}

	if strings.Contains(houseType, "woning") {
		return corporation.HousingTypeHouse
	}

	return corporation.HousingTypeUndefined
}
