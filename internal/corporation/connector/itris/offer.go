package itris

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/woningfinder/woningfinder/internal/corporation"
	"go.uber.org/zap"
)

const detailsHousingChildAttr = "li.link a"

func (c *client) FetchOffers(ch chan<- corporation.Offer) error {
	offers := map[string]*corporation.Offer{}

	// create another collector for housing details
	detailCollector := c.collector.Clone()

	// add offer
	c.collector.OnHTML("div.aanbodListItems", func(el *colly.HTMLElement) {
		el.ForEach("div.woningaanbod", func(_ int, e *colly.HTMLElement) {
			var offer corporation.Offer
			var err error

			// set corporation name
			offer.CorporationName = c.corporation.Name

			// get offer url
			offer.URL = c.corporation.APIEndpoint.String() + e.ChildAttr(detailsHousingChildAttr, "href")
			offer.ExternalID = e.Attr("data-aanbod-id")

			// get picture url
			offer.RawPictureURL, err = c.parsePictureURL(e.ChildAttr("div.image-wrapper", "style"))
			if err != nil {
				c.logger.Info("failed parsing picture url", zap.Error(err), logConnector)
			}

			// get housing type
			offer.Housing.Type = c.parseHousingType(e.Text)
			if offer.Housing.Type == corporation.HousingTypeUndefined {
				return
			}

			// get location
			offer.Housing.Address = strings.Title(strings.ToLower(e.ChildAttr(detailsHousingChildAttr, "data-select-address")))
			offer.Housing.CityName = strings.Title(strings.ToLower(e.Attr("data-plaats")))
			offer.Housing.CityDistrict, err = c.mapboxClient.CityDistrictFromAddress(offer.Housing.Address)
			if err != nil {
				c.logger.Info("could not get city district", zap.String("address", offer.Housing.Address), zap.Error(err), logConnector)
			}

			offer.Housing.Price, err = strconv.ParseFloat(e.Attr("data-prijs"), 32)
			if err != nil {
				c.logger.Info("error while parsing price", zap.String("address", offer.Housing.Address), zap.Error(err), logConnector)
			}

			offer.Housing.NumberBedroom, err = strconv.Atoi(e.Attr("data-kamers"))
			if err != nil {
				c.logger.Info("error while parsing number bedroom", zap.String("address", offer.Housing.Address), zap.Error(err), logConnector)
			}

			// create new offer
			offers[offer.URL] = &offer

			// visit offer url
			if err := detailCollector.Visit(offer.URL); err != nil {
				c.logger.Warn("error while checking offer details", zap.String("address", offer.Housing.Address), zap.Error(err), logConnector)
			}
		})
	})

	// add housing details
	detailCollector.OnHTML("html", func(e *colly.HTMLElement) {
		// find offer
		offerURL := e.Request.URL.String()
		if _, ok := offers[offerURL]; !ok {
			// no offer matching, no details
			return
		}

		if err := c.getHousingDetails(offers[offerURL], e); err != nil {
			c.logger.Info("error while getting house details", zap.String("address", offers[offerURL].Housing.Address), zap.Error(err), logConnector)
			return
		}

		ch <- *offers[offerURL]
	})

	// parse offers
	offerURL := c.corporation.APIEndpoint.String() + "/woningaanbod/"
	if err := c.collector.Visit(offerURL); err != nil {
		return err
	}

	return nil
}

func (c *client) getHousingDetails(offer *corporation.Offer, e *colly.HTMLElement) error {
	// add housing size
	e.ForEach("div.info-container > #oppervlaktes-page div.infor-wrapper", func(_ int, el *colly.HTMLElement) {
		// increase size
		roomSize, err := strconv.ParseFloat(strings.ReplaceAll(strings.Trim(el.Text, " m2"), ",", "."), 64)
		if err != nil {
			return
		}
		offer.Housing.Size += roomSize
	})

	// add housing details
	e.ForEach("#icons-container ul > li", func(_ int, el *colly.HTMLElement) {
		switch el.Text {
		case "Balkon":
			offer.Housing.Balcony = el.DOM.HasClass("yes")
		case "Garage":
			offer.Housing.Garage = el.DOM.HasClass("yes")
		case "Lift":
			offer.Housing.Elevator = el.DOM.HasClass("yes")
		case "Tuin":
			offer.Housing.Garden = el.DOM.HasClass("yes")
		}
	})

	// part of housing details can be found in housing description
	offer.Housing.Accessible = strings.Contains(e.Text, "toegankelijk")

	return nil
}

func (c *client) parseHousingType(houseType string) corporation.HousingType {
	houseType = strings.ToLower(houseType)

	if strings.Contains(houseType, "appartement") ||
		strings.Contains(houseType, "penthouse") {
		return corporation.HousingTypeAppartement
	}

	if strings.Contains(houseType, "woning") {
		return corporation.HousingTypeHouse
	}

	return corporation.HousingTypeUndefined
}

func (c *client) parsePictureURL(path string) (*url.URL, error) {
	if path == "" {
		return nil, nil
	}

	// format path
	path = strings.ReplaceAll(path, "background-image:url(", "")
	path = strings.ReplaceAll(path, ");", "")

	pictureURL, err := url.Parse(c.corporation.APIEndpoint.String() + path)
	if err != nil {
		return nil, fmt.Errorf("failed to parse picture url %s: %w", path, err)
	}

	return pictureURL, nil
}
