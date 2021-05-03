package onshuis

import (
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/woningfinder/woningfinder/internal/corporation"
	"github.com/woningfinder/woningfinder/pkg/logging"
)

func DetailsParser(logger *logging.Logger, offer *corporation.Offer, e *colly.HTMLElement) {
	// add selection method
	offer.SelectionMethod = corporation.SelectionRandom // all houses from onshuis are random

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

	dom, err := e.DOM.Html()
	if err != nil {
		logger.Sugar().Warnf("unable to get details for %s on %s", offer.Housing.Address, offer.URL)
		return
	}

	// add garden
	offer.Housing.Garage = strings.Contains(dom, "tuin")

	// add garage
	offer.Housing.Garage = strings.Contains(dom, "garage")

	// add elevator
	offer.Housing.Elevator = strings.Contains(dom, "lift")

	// add balcony
	offer.Housing.Balcony = strings.Contains(dom, "balkon")

	// add attic
	offer.Housing.Attic = strings.Contains(dom, "zolder")

	// add accessible
	offer.Housing.Accessible = strings.Contains(dom, "toegankelijk")
}
