package itris

import (
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/woningfinder/woningfinder/internal/corporation"
)

const detailsHousingChildAttr = "li.link a"

func (c *itrisConnector) FetchOffer() ([]corporation.Offer, error) {
	offerURL := c.url + "/woningaanbod/"

	//get offer
	var offers []corporation.Offer
	c.collector.OnHTML("div.aanbodListItems", func(el *colly.HTMLElement) {
		el.ForEach("div.woningaanbod", func(_ int, e *colly.HTMLElement) {
			houseType := c.parseHousingType(e.Text)
			if houseType == corporation.Undefined {
				return
			}

			address := strings.Title(strings.ToLower(e.ChildAttr(detailsHousingChildAttr, "data-select-address")))
			latitude, longitude := c.parseLocation(e.ChildAttr(detailsHousingChildAttr, "data-select-lat-long"), address)
			if latitude == 0 || longitude == 0 {
				log.Printf("error while parsing coordinates of %s: [%f, %f]", address, latitude, longitude)
				return
			}

			reactionDate, err := time.Parse(layoutTime, e.Attr("data-reactiedatum"))
			if err != nil {
				log.Printf("error while parsing date of %s: %v", address, err)
				return
			}

			price, err := strconv.ParseFloat(e.Attr("data-prijs"), 32)
			if err != nil {
				log.Printf("error while parsing price of %s: %v", address, err)
				return
			}

			numberBedroom, err := strconv.Atoi(e.Attr("data-kamers"))
			if err != nil {
				log.Printf("error while parsing number bedroom of %s: %v", address, err)
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

			offers = append(offers, offer)
		})
	})

	if err := c.collector.Visit(offerURL); err != nil {
		return nil, err
	}

	return offers, nil
}

func (c *itrisConnector) parseHousingType(houseType string) corporation.Type {
	houseType = strings.ToLower(houseType)

	if strings.Contains(houseType, "appartement") {
		return corporation.Appartement
	}

	if strings.Contains(houseType, "eengezinswoning") {
		return corporation.House
	}

	return corporation.Undefined
}

func (c *itrisConnector) parseLocation(entry, address string) (latitude float64, longitude float64) {
	// latitude is idx 0 and longitude is 1
	coordinates := strings.Split(entry, "-")
	if len(coordinates) != 2 {
		log.Printf("error while parsing coordinates of %s: %v", address, coordinates)
		return
	}
	latitude, err := strconv.ParseFloat(coordinates[0], 32)
	if err != nil {
		log.Printf("error while parsing latitude of %s: %v", address, err)
		return
	}

	longitude, err = strconv.ParseFloat(coordinates[1], 32)
	if err != nil {
		log.Printf("error while parsing latitude of %s: %v", address, err)
		return
	}

	return
}
