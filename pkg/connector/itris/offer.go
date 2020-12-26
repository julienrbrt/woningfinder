package itris

import (
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/woningfinder/woningfinder/internal/corporation"
)

const detailsHousingChildAttr = "li.link a"

func (c *itrisConnector) FetchOffer() ([]corporation.Offer, error) {
	offers := map[string]*corporation.Offer{}

	// add offer
	c.collector.OnHTML("div.aanbodListItems", func(el *colly.HTMLElement) {
		el.ForEach("div.woningaanbod", func(_ int, e *colly.HTMLElement) {
			houseType := c.parseHousingType(e.Text)
			if houseType == corporation.Undefined {
				return
			}

			address := strings.Title(strings.ToLower(e.ChildAttr(detailsHousingChildAttr, "data-select-address")))
			latitude, longitude := c.parseLocation(e.ChildAttr(detailsHousingChildAttr, "data-select-lat-long"), address)
			if latitude == 0 || longitude == 0 {
				c.logger.Sugar().Errorf("error while parsing coordinates of %s: [%f, %f]", address, latitude, longitude)
				return
			}

			reactionDate, err := time.Parse(layoutTime, e.Attr("data-reactiedatum"))
			if err != nil {
				c.logger.Sugar().Errorf("error while parsing date of %s: %w", address, err)
				return
			}

			price, err := strconv.ParseFloat(e.Attr("data-prijs"), 16)
			if err != nil {
				c.logger.Sugar().Errorf("error while parsing price of %s: %w", address, err)
				return
			}

			numberBedroom, err := strconv.Atoi(e.Attr("data-kamers"))
			if err != nil {
				c.logger.Sugar().Errorf("error while parsing number bedroom of %s: %w", address, err)
				return
			}

			offer := corporation.Offer{
				SelectionDate: reactionDate,
				URL:           c.url + e.ChildAttr(detailsHousingChildAttr, "href"),
				ExternalID:    e.Attr("data-aanbod-id"),
				Housing: corporation.Housing{
					Type: corporation.HousingType{
						Type: houseType,
					},
					Price:   price,
					Address: address,
					City: corporation.City{
						Name: strings.Title(strings.ToLower(e.Attr("data-plaats"))),
					},
					Latitude:      latitude,
					Longitude:     longitude,
					NumberBedroom: numberBedroom,
				},
			}

			// create new offer
			offers[offer.URL] = &offer

			// visit offer url
			if err := e.Request.Visit(offer.URL); err != nil {
				c.logger.Sugar().Errorf("error while checking offer details %s: %w", address, err)
			}
		})
	})

	// add housing details
	c.collector.OnHTML("div.info-container", func(e *colly.HTMLElement) {
		// Find offer
		offerURL := e.Request.URL.String()
		if _, ok := offers[offerURL]; !ok {
			// no offer matching, no details
			return
		}

		// add number of room
		var numberRoom int
		e.ForEach("#oppervlaktes-page h3", func(_ int, _ *colly.HTMLElement) {
			numberRoom++
		})
		if numberRoom > 0 {
			offers[offerURL].Housing.NumberRoom = numberRoom
		}

		// add energie label
		energieLabel := e.ChildText("#Woning-page strong.tag-text")
		if energieLabel != "" {
			offers[offerURL].Housing.EnergieLabel = energieLabel
		}

		// add building year
		e.ForEach("div.infor-wrapper", func(_ int, el *colly.HTMLElement) {
			buildingYear, err := strconv.Atoi(el.Text)
			if err != nil {
				return
			}
			if buildingYear > 1850 { // random building year so high that it cannot be a number of room
				offers[offerURL].Housing.BuildingYear = buildingYear
			}
		})
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

func (c *itrisConnector) parseHousingType(houseType string) corporation.Type {
	houseType = strings.ToLower(houseType)

	if strings.Contains(houseType, "appartement") {
		return corporation.Appartement
	}

	if strings.Contains(houseType, "woning") {
		return corporation.House
	}

	return corporation.Undefined
}

func (c *itrisConnector) parseLocation(entry, address string) (latitude float64, longitude float64) {
	// latitude is idx 0 and longitude is 1
	coordinates := strings.Split(entry, "-")
	if len(coordinates) != 2 {
		c.logger.Sugar().Warnf("error while parsing coordinates of %s: %v", address, coordinates)
		return
	}

	latitude, err := strconv.ParseFloat(coordinates[0], 16)
	if err != nil {
		c.logger.Sugar().Errorf("error while parsing latitude of %s: %v", address, err)
		return
	}

	longitude, err = strconv.ParseFloat(coordinates[1], 16)
	if err != nil {
		c.logger.Sugar().Errorf("error while parsing latitude of %s: %v", address, err)
		return
	}

	return
}
