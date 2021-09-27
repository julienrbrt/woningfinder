package domijn

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/woningfinder/woningfinder/internal/corporation"
)

func (c *client) GetOffers() ([]corporation.Offer, error) {
	offers := map[string]*corporation.Offer{}

	// create another collector for housing details
	paginationCollector := c.collector.Clone()
	detailCollector := c.collector.Clone()

	// check paginating
	c.collector.OnHTML("ol[aria-labelledby=ARIA-Label-paging]", func(el *colly.HTMLElement) {
		el.ForEach("li", func(_ int, e *colly.HTMLElement) {
			// visit other page
			paginatedURL := e.ChildAttr("a", "href")
			if paginatedURL == "" { // first page
				return
			}

			if err := paginationCollector.Visit(c.corporation.APIEndpoint.String() + paginatedURL); err != nil {
				c.logger.Sugar().Warnf("domijn connector: error while checking pagination %s: %w", c.corporation.APIEndpoint.String()+paginatedURL, err)
			}
		})
	})

	offerParser := func(el *colly.HTMLElement) {
		el.ForEach("article", func(_ int, e *colly.HTMLElement) {
			var offer corporation.Offer
			var err error

			// set selection method - domijn is always random selection
			offer.SelectionMethod = corporation.SelectionRandom

			// get offer url
			offer.URL = c.corporation.APIEndpoint.String() + strings.ReplaceAll(e.ChildAttr("div.image-wrapper a", "href"), " ", "%20") // manual space escaping

			// get image url
			offer.RawPictureURL, err = c.parsePictureURL(e.ChildAttr("div.image-wrapper a img", "src"))
			if err != nil {
				c.logger.Sugar().Info(err)
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
				c.logger.Sugar().Infof("domijn connector: could not get city district of %s: %w", offer.Housing.Address, err)
			}

			// create new offer
			offers[offer.URL] = &offer

			// visit offer url
			if err := detailCollector.Visit(offer.URL); err != nil {
				c.logger.Sugar().Warnf("domijn connector: error while checking offer details %s: %w", offer.Housing.Address, err)
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

		c.getHousingDetails(offers[offerURL], e)
	})

	// parse offers
	offerURL := c.corporation.APIEndpoint.String() + "/woningaanbod/?advertisementType=rent"
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

func (c *client) getHousingDetails(offer *corporation.Offer, e *colly.HTMLElement) {
	var err error

	// parse externalID
	offer.ExternalID = e.ChildAttr("form[name=reactionform]", "action")

	// parse price
	priceStr := strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(e.ChildText("span.price"), ",", "."), "€", ""))
	offer.Housing.Price, err = strconv.ParseFloat(priceStr, 32)
	if err != nil {
		c.logger.Sugar().Infof("domijn connector: error while parsing price of %s: %w", offer.Housing.Address, err)
		return
	}

	// parse housing characteristics
	e.ForEach("div.properties span > span", func(_ int, el *colly.HTMLElement) {
		property := strings.TrimSpace(el.ChildText("strong"))

		switch property {
		case "Energielabel":
			offer.Housing.EnergyLabel = el.ChildText("span.c-energylabel")
		case "Slaapkamers":
			offer.Housing.NumberBedroom, err = strconv.Atoi(cleanProperty(el.Text, property))
			if err != nil {
				return
			}
		case "Tuin / Balkon":
			switch strings.ToLower(cleanProperty(el.Text, property)) {
			case "balkon":
				offer.Housing.Balcony = true
			case "tuin":
				offer.Housing.Garden = true
			}
		case "Bouwjaar":
			offer.Housing.BuildingYear, err = strconv.Atoi(cleanProperty(el.Text, property))
			if err != nil {
				return
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
			offer.MinFamilySize = 4
		}
	})
}

func (c *client) parseHousingType(houseType string) corporation.HousingType {
	houseType = strings.ToLower(houseType)

	if strings.Contains(houseType, "appartement") {
		return corporation.HousingTypeAppartement
	}

	if strings.Contains(houseType, "woning") || strings.Contains(houseType, "duplex") {
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
		return nil, fmt.Errorf("domijn connector: failed to parse picture url %s: %w", path, err)
	}

	return pictureURL, nil
}

func cleanProperty(input, property string) string {
	return strings.TrimSpace(strings.ReplaceAll(input, property, ""))
}
