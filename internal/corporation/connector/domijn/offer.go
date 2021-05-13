package domijn

import (
	"fmt"
	"strconv"
	"strings"
	"time"

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

			if err := paginationCollector.Visit(c.url + paginatedURL); err != nil {
				c.logger.Sugar().Errorf("domijn connector: error while checking pagination %s: %w", c.url+paginatedURL, err)
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
			offer.URL = c.url + strings.ReplaceAll(e.ChildAttr("div.image-wrapper a", "href"), " ", "%20") // manual space escaping

			// get housing type
			offer.Housing.Type = c.parseHousingType(e.ChildText("div.info div.row div p"))
			if offer.Housing.Type == corporation.HousingTypeUndefined {
				return
			}

			// get location
			offer.Housing.City.Name = strings.TrimSpace(e.ChildText("div.info > p"))
			offer.Housing.Address = fmt.Sprintf("%s, %s", e.ChildText("div.info > h1"), offer.Housing.City.Name)
			offer.Housing.CityDistrict, err = c.mapboxClient.CityDistrictFromAddress(offer.Housing.Address)
			if err != nil {
				c.logger.Sugar().Warnf("domijn connector: could not get city district of %s: %w", offer.Housing.Address, err)
			}

			// create new offer
			offers[offer.URL] = &offer

			// visit offer url
			if err := detailCollector.Visit(offer.URL); err != nil {
				c.logger.Sugar().Errorf("domijn connector: error while checking offer details %s: %w", offer.Housing.Address, err)
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

func (c *client) getHousingDetails(offer *corporation.Offer, e *colly.HTMLElement) {
	var err error

	// parse externalID
	offer.ExternalID = e.ChildAttr("form[name=reactionform]", "action")

	// parse price
	priceRaw := strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(e.ChildText("span.price"), ",", "."), "€", ""))
	offer.Housing.Price, err = strconv.ParseFloat(priceRaw, 16)
	if err != nil {
		c.logger.Sugar().Errorf("domijn connector: error while parsing price of %s: %w", offer.Housing.Address, err)
		return
	}

	// parse selection date
	e.ForEach("div.card dd", func(i int, el *colly.HTMLElement) {
		// the date is at the 3rd position in the table
		if i == 2 {
			offer.SelectionDate, err = time.Parse(layoutTime, strings.TrimSpace(el.Text))
			if err != nil {
				c.logger.Sugar().Errorf("domijn connector: error while parsing date of %s: %w", offer.Housing.Address, err)
				return
			}
		}
	})

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

func cleanProperty(input, property string) string {
	return strings.TrimSpace(strings.ReplaceAll(input, property, ""))
}
