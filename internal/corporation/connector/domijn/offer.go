package domijn

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/woningfinder/woningfinder/internal/corporation"
	"go.uber.org/zap"
)

func (c *client) FetchOffers(ch chan<- corporation.Offer) error {
	offers := map[string]*corporation.Offer{}

	// create another collector for housing details
	paginationCollector := c.collector.Clone()
	detailCollector := c.collector.Clone()

	// check paginating
	c.collector.OnHTML("ol[aria-labelledby=ARIA-Label-paging]", func(el *colly.HTMLElement) {
		el.ForEach("li", func(_ int, e *colly.HTMLElement) {
			// visit other pages
			paginatedURL := e.ChildAttr("a", "href")
			if paginatedURL == "" { // first page
				return
			}

			if err := paginationCollector.Visit(c.corporation.APIEndpoint.String() + paginatedURL); err != nil {
				c.logger.Warn("error while checking pagination", zap.String("url", c.corporation.APIEndpoint.String()+paginatedURL), zap.Error(err), logConnector)
			}
		})
	})

	offerParser := func(el *colly.HTMLElement) {
		el.ForEach("article", func(_ int, e *colly.HTMLElement) {
			var offer corporation.Offer
			var err error

			// set corporation name
			offer.CorporationName = c.corporation.Name

			// get offer url
			offer.URL = c.corporation.APIEndpoint.String() + strings.ReplaceAll(e.ChildAttr("div.image-wrapper a", "href"), " ", "%20") // manual space escaping

			// get image url
			offer.RawPictureURL, err = c.parsePictureURL(e.ChildAttr("div.image-wrapper a img", "src"))
			if err != nil {
				c.logger.Info("failed parsing picture url", zap.Error(err), logConnector)
			}

			// get housing type
			offer.Housing.Type = c.parseHousingType(e.ChildText("div.info div.row div p"))
			if offer.Housing.Type == corporation.HousingTypeUndefined {
				return
			}

			// get location
			offer.Housing.CityName = strings.TrimSpace(e.ChildText("div.info > p"))
			offer.Housing.Address = fmt.Sprintf("%s, %s", e.ChildText("div.info > h1"), offer.Housing.CityName)
			offer.Housing.CityDistrict, err = c.mapboxClient.CityDistrictFromAddress(offer.Housing.Address)
			if err != nil {
				c.logger.Info("could not get city district", zap.String("address", offer.Housing.Address), zap.Error(err), logConnector)
			}

			// create new offer
			offers[offer.URL] = &offer

			// visit offer url
			if err := detailCollector.Visit(offer.URL); err != nil {
				c.logger.Warn("error while checking offer details", zap.String("address", offer.Housing.Address), zap.Error(err), logConnector)
			}
		})
	}

	// add offer
	c.collector.OnHTML("#housingWrapper", offerParser) // parses first page
	paginationCollector.OnHTML("#housingWrapper", offerParser)

	// add housing details
	detailCollector.OnHTML("#content-anchor", func(e *colly.HTMLElement) {
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
	offerURL := c.corporation.APIEndpoint.String() + "/woningaanbod/?advertisementType=rent"
	if err := c.collector.Visit(offerURL); err != nil {
		return err
	}

	return nil
}

func (c *client) getHousingDetails(offer *corporation.Offer, e *colly.HTMLElement) error {
	var err error

	// parse externalID
	offer.ExternalID = e.ChildAttr("form[name=reactionform]", "action")

	// parse price
	priceStr := strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(e.ChildText("span.price"), ",", "."), "€", ""))
	offer.Housing.Price, err = strconv.ParseFloat(priceStr, 32)
	if err != nil {
		return fmt.Errorf("error while parsing price: %w", err)
	}

	// parse housing characteristics
	e.ForEach("div.properties span > span", func(_ int, el *colly.HTMLElement) {
		property := strings.TrimSpace(el.ChildText("strong"))

		switch property {
		case "Slaapkamers":
			offer.Housing.NumberBedroom, err = strconv.Atoi(cleanProperty(el.Text, property))
			if err != nil {
				c.logger.Info("error parsing number bedroom", zap.String("address", offer.Housing.Address), zap.Error(err), logConnector)
			}
		case "Tuin / Balkon":
			switch strings.ToLower(cleanProperty(el.Text, property)) {
			case "balkon":
				offer.Housing.Balcony = true
			case "tuin":
				offer.Housing.Garden = true
			}
		case "Lift":
			offer.Housing.Elevator = strings.EqualFold(cleanProperty(el.Text, property), "ja")
		case "Seniorenwoning":
			// from domijn seniorenwoning are from 65+
			if strings.EqualFold(cleanProperty(el.Text, property), "ja") {
				offer.MinAge = 65
			}
		}

		// logic for minimum family size https://www.domijn.nl/ik-zoek/huurwoningen/huren-grote-eengezinswoning/
		if offer.Housing.NumberBedroom >= 4 {
			offer.MinFamilySize = 3
		}
	})

	return nil
}

func (c *client) parseHousingType(houseType string) corporation.HousingType {
	houseType = strings.ToLower(houseType)

	if strings.Contains(houseType, "appartement") {
		return corporation.HousingTypeAppartement
	}

	if strings.Contains(houseType, "woning") ||
		strings.Contains(houseType, "duplex") {
		return corporation.HousingTypeHouse
	}

	return corporation.HousingTypeUndefined
}

func (c *client) parsePictureURL(path string) (*url.URL, error) {
	if path == "" {
		return nil, nil
	}

	pictureURL, err := url.Parse(c.corporation.APIEndpoint.String() + strings.ReplaceAll(path, "280/190", "600/400"))
	if err != nil {
		return nil, fmt.Errorf("failed to parse picture url %s: %w", path, err)
	}

	return pictureURL, nil
}

func cleanProperty(input, property string) string {
	return strings.TrimSpace(strings.ReplaceAll(input, property, ""))
}
