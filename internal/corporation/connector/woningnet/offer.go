package woningnet

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
	"github.com/woningfinder/woningfinder/internal/corporation"
	"github.com/woningfinder/woningfinder/pkg/networking"
	"go.uber.org/zap"
)

type response struct {
	Offer   []offer `json:"Resultaten"`
	Filters struct {
		FiltersList []struct {
			GetTypeString string `json:"GetTypeString"`
			Paging        struct {
				CurrentPage    int `json:"CurrentPage"`
				TotalPageCount int `json:"TotalPageCount"`
			} `json:"Paging,omitempty"`
		} `json:"FiltersList"`
	} `json:"Filters"`
}

type offer struct {
	Address                        string    `json:"Adres"`
	CityAndDistrict                string    `json:"PlaatsWijk"`
	Prijs                          string    `json:"Prijs"`
	RawPictureURL                  string    `json:"AfbeeldingUrl"`
	SelectionDate                  time.Time `json:"PublicatieEinddatum"`
	PublicatieBegindatum           time.Time `json:"PublicatieBegindatum"`
	PublicatieID                   string    `json:"PublicatieId"`
	AdvertentieURL                 string    `json:"AdvertentieUrl"`
	PublicatieModel                string    `json:"PublicatieModel"`
	CurrentRegioCode               string    `json:"CurrentRegioCode"`
	Woonoppervlakte                string    `json:"Woonoppervlakte"`
	HouseType                      string    `json:"DetailSoortOmschrijving"`
	ToegankelijkheidslabelCSSClass string    `json:"ToegankelijkheidslabelCssClass"`
}

// offerRequest builds a WoningNet request
func offerRequest(commandURL, command string) (networking.Request, error) {
	req := struct {
		URL       string `json:"url"`
		Command   string `json:"command"`
		Hideunits string `json:"hideunits"`
	}{
		URL:       commandURL,
		Command:   command,
		Hideunits: "hideunits[]",
	}

	body, err := json.Marshal(req)
	if err != nil {
		return networking.Request{}, fmt.Errorf("error while marshaling %v: %w", req, err)
	}

	request := networking.Request{
		Path:   "/zoeken/find/",
		Method: http.MethodPost,
		Body:   bytes.NewBuffer(body),
	}

	return request, nil
}

// selectionMethodToCommand map the selection methods to woningnet requests commands
var selectionMethodCommandMap = map[corporation.SelectionMethod]string{
	corporation.SelectionRegistrationDate:     "model[Regulier%20aanbod]",
	corporation.SelectionFirstComeFirstServed: "model[Vrijesectorhuur]",
	corporation.SelectionRandom:               "model[Loting]",
}

func (c *client) FetchOffers(ch chan<- corporation.Offer) error {
	for _, selectionMethod := range c.corporation.SelectionMethod {
		req, err := offerRequest(selectionMethodCommandMap[selectionMethod], "")
		if err != nil {
			return err
		}

		resp, err := c.Send(req)
		if err != nil {
			return err
		}

		var response response
		if err := json.Unmarshal(resp, &response); err != nil {
			return fmt.Errorf("error parsing offer response %v: %w", string(resp), err)
		}

		paginatedOffers, err := c.getPaginatedOffers(selectionMethodCommandMap[selectionMethod], response)
		if err != nil {
			c.logger.Warn("error getting paginated offers", zap.Error(err), logConnector)
		}
		respOffers := append(response.Offer, paginatedOffers...)

		for _, offer := range respOffers {
			houseType := c.parseHousingType(offer.HouseType)
			if houseType == corporation.HousingTypeUndefined {
				continue
			}

			result, err := c.Map(offer, houseType)
			if err != nil {
				c.logger.Warn("error mapping woningnet offer", zap.String("address", offer.Address), zap.Error(err), logConnector)
				continue
			}

			// add offer to channel
			ch <- result
		}
	}

	return nil
}

// get PaginatedOffers get the wonignnet offers that are paginated
func (c *client) getPaginatedOffers(commandURL string, resp response) ([]offer, error) {
	var offers []offer

	for _, f := range resp.Filters.FiltersList {
		if f.GetTypeString != "PagingFilter" {
			continue
		}

		for i := f.Paging.TotalPageCount; i > f.Paging.CurrentPage; i-- {
			req, err := offerRequest(commandURL, fmt.Sprintf("page[%d]", i))
			if err != nil {
				return nil, err
			}

			resp, err := c.Send(req)
			if err != nil {
				return nil, err
			}

			var newResp response
			if err := json.Unmarshal(resp, &newResp); err != nil {
				return nil, fmt.Errorf("error parsing offer response %v: %w", string(resp), err)
			}

			offers = append(offers, newResp.Offer...)
		}
	}

	return offers, nil
}

func (c *client) Map(offer offer, houseType corporation.HousingType) (corporation.Offer, error) {
	var err error
	var offerURL = c.corporation.URL + offer.AdvertentieURL
	if strings.Contains(offer.CityAndDistrict, "%") { // escaprd character in city name
		if offer.CityAndDistrict, err = url.QueryUnescape(offer.CityAndDistrict); err != nil {
			c.logger.Warn("failed unescaping city", zap.String("address", offer.Address), zap.Error(err), logConnector)
		}
	}
	var cityName = strings.Split(offer.CityAndDistrict, " - ")

	// fill in already known housing characteristics
	house := corporation.Housing{
		Type:       houseType,
		Address:    fmt.Sprintf("%s %s", offer.Address, cityName[0]),
		CityName:   cityName[0],
		Accessible: offer.ToegankelijkheidslabelCSSClass != "ToegankelijkheidslabelZON",
	}

	house.Price, err = c.parsePrice(offer.Prijs)
	if err != nil {
		return corporation.Offer{}, fmt.Errorf("error mapping offer %s: %w", offerURL, err)
	}

	house.Size, err = c.parseHouseSize(offer.Woonoppervlakte)
	if err != nil {
		c.logger.Info("failed parsing house size", zap.String("address", house.Address), zap.Error(err), logConnector)
	}

	// get address city district
	house.CityDistrict, err = c.mapboxClient.CityDistrictFromAddress(house.Address)
	if err != nil {
		house.CityDistrict = cityName[len(cityName)-1]
		c.logger.Info("could not get city district", zap.String("address", house.Address), zap.Error(err), logConnector)
	}

	// get housing picture url
	rawPictureURL, err := c.parsePictureURL(offer.RawPictureURL)
	if err != nil {
		c.logger.Info("failed parsing picture url", zap.Error(err), logConnector)
	}

	c.collector.OnHTML("#Kenmerken", func(e *colly.HTMLElement) {
		table := e.DOM.ChildrenFiltered(".contentBlock")

		// number bedroom
		house.NumberBedroom, err = c.parseBedroom(c.getContentValue("Aantal kamers (incl. woonkamer)", table))
		if err != nil {
			c.logger.Info("error parsing number bedroom", zap.String("address", house.Address), zap.Error(err), logConnector)
		}

		// outside
		outside := c.getContentValue("Buitenruimte", table)
		house.Balcony = strings.Contains(outside, "balkon") || strings.Contains(outside, "terras")
		house.Garden = strings.Contains(outside, "tuin")
		house.Garage = strings.Contains(outside, "garage") || strings.Contains(outside, "parkeer")
		house.Elevator = strings.Contains(e.Text, "lift")
	})

	// parse offer conditions
	var minAge int
	c.collector.OnHTML("div.publicationContainer", func(e *colly.HTMLElement) {
		if strings.Contains(e.Text, "Senioren") {
			minAge = 55
		}
	})

	if err = c.collector.Visit(offerURL); err != nil {
		return corporation.Offer{}, fmt.Errorf("error visiting %s: %w", offerURL, err)
	}

	return corporation.Offer{
		CorporationName: c.corporation.Name,
		ExternalID:      fmt.Sprintf("%s/%s", offer.PublicatieID, offer.CurrentRegioCode),
		Housing:         house,
		URL:             offerURL,
		RawPictureURL:   rawPictureURL,
		MinAge:          minAge,
	}, nil
}

func (c *client) parseHousingType(houseType string) corporation.HousingType {
	houseType = strings.ToLower(houseType)

	if strings.EqualFold(houseType, "Galerijflat") ||
		strings.EqualFold(houseType, "Portiekflat") ||
		strings.EqualFold(houseType, "Benedenwoning") ||
		strings.EqualFold(houseType, "Bovenwoning") ||
		strings.EqualFold(houseType, "Corridorflat") ||
		strings.EqualFold(houseType, "Maisonnette") ||
		strings.Contains(houseType, "flat") {
		return corporation.HousingTypeAppartement
	}

	if strings.Contains(houseType, "woning") {
		return corporation.HousingTypeHouse
	}

	return corporation.HousingTypeUndefined
}

func (c *client) parsePrice(priceStr string) (float64, error) {
	if strings.Contains(priceStr, " - ") {
		priceStr = strings.Split(priceStr, " - ")[0]
	}

	// remove dot delimiter in price
	priceStr = strings.ReplaceAll(priceStr, ".", "")

	price, err := strconv.ParseFloat(strings.ReplaceAll(strings.TrimLeft(priceStr, "€ "), ",", "."), 64)
	if err != nil {
		return 0, fmt.Errorf("error parsing housing price %s: %w", priceStr, err)
	}

	return price, nil
}

func (c *client) parseHouseSize(sizeStr string) (float64, error) {
	if strings.Contains(sizeStr, " - ") {
		sizeStr = strings.Split(sizeStr, " - ")[0]
	}

	size, err := strconv.ParseFloat(sizeStr, 64)
	if err != nil {
		return 0, fmt.Errorf("error parsing housing size %s: %w", sizeStr, err)
	}

	return size, nil
}

func (c *client) parseBedroom(bedroomStr string) (int, error) {
	var err error
	var bedroom int

	if strings.Contains(bedroomStr, "slaap") {
		bedroomStr = string([]rune(strings.Split(bedroomStr, "(")[1])[0])
		bedroom, err = strconv.Atoi(bedroomStr)
		if err != nil {
			return bedroom, fmt.Errorf("error parsing number bedrooms %s: %w", bedroomStr, err)
		}
	}

	return bedroom, nil
}

// getContentValue is use to get the value from a table in the woningnet housing description page
func (c *client) getContentValue(name string, selection *goquery.Selection) string {
	var value string
	selection.Children().Each(func(i int, s *goquery.Selection) {
		if strings.EqualFold(s.Children().First().Text(), name) {
			value = s.Children().Next().Text()
		}
	})

	return strings.ToLower(value)
}

func (c *client) parsePictureURL(path string) (*url.URL, error) {
	if path == "" {
		return nil, nil
	}

	pictureURL, err := url.Parse(path)
	if err != nil {
		return nil, fmt.Errorf("failed to parse picture url %s: %w", path, err)
	}

	return pictureURL, nil
}
