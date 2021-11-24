package zig

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"

	"github.com/julienrbrt/woningfinder/internal/corporation"
	"github.com/julienrbrt/woningfinder/pkg/networking"
	"go.uber.org/zap"
)

const externalIDSeperator = ";"

type offerList struct {
	Dwellingtype struct {
		Categorie     string `json:"categorie"`
		Localizedname string `json:"localizedName"`
	} `json:"dwellingType"`
	TotalRent float64 `json:"totalRent"`
	Rentbuy   string  `json:"rentBuy"`
	ID        string  `json:"id"`
}

type offerDetails struct {
	Postalcode          string `json:"postalcode"`
	Street              string `json:"street"`
	Housenumber         string `json:"houseNumber"`
	Housenumberaddition string `json:"houseNumberAddition"`
	City                struct {
		Name string `json:"name"`
	} `json:"city"`
	Dwellingtype struct {
		Categorie     string `json:"categorie"`
		Localizedname string `json:"localizedName"`
	} `json:"dwellingType"`
	Totalrent    float64 `json:"totalRent"`
	Sleepingroom struct {
		Amountofrooms string `json:"amountOfRooms"`
		ID            string `json:"id"`
		Localizedname string `json:"localizedName"`
	} `json:"sleepingRoom"`
	Garden               bool   `json:"garden"`
	Balcony              bool   `json:"balcony"`
	Minimumincome        int    `json:"minimumIncome"`
	Maximumincome        int    `json:"maximumIncome"`
	Minimumhouseholdsize int    `json:"minimumHouseholdSize"`
	Maximumhouseholdsize int    `json:"maximumHouseholdSize"`
	Minimumage           int    `json:"minimumAge"`
	Maximumage           int    `json:"maximumAge"`
	Rentbuy              string `json:"rentBuy"`
	Assignmentid         int    `json:"assignmentID"`
	Pictures             []struct {
		Label string `json:"label"`
		URI   string `json:"uri"`
		Type  string `json:"type"`
	} `json:"pictures"`
	Areadwelling  float64 `json:"areaDwelling"`
	Iszelfstandig bool    `json:"isZelfstandig"`
	Urlkey        string  `json:"urlKey"`
	ID            string  `json:"id"`
}

func offerRequest() networking.Request {
	body := url.Values{}
	body.Add("configurationKeys[]", "aantalReacties")
	body.Add("configurationKeys[]", "passend")

	request := networking.Request{
		Path:   "/portal/object/frontend/getallobjects/format/json",
		Method: http.MethodPost,
		Body:   strings.NewReader(body.Encode()),
	}

	return request
}

func offerDetailRequest(offerID string) networking.Request {
	body := url.Values{}
	body.Add("id", offerID)

	request := networking.Request{
		Path:   "/portal/object/frontend/getobject/format/json",
		Method: http.MethodPost,
		Body:   strings.NewReader(body.Encode()),
	}

	return request
}

func (c *client) FetchOffers(ch chan<- corporation.Offer) error {
	resp, err := c.Send(offerRequest())
	if err != nil {
		return err
	}

	var result struct {
		Result []offerList `json:"result"`
	}
	if err := json.Unmarshal(resp, &result); err != nil {
		return fmt.Errorf("error parsing offer result %v: %w", string(resp), err)
	}

	var wg sync.WaitGroup
	for _, offer := range result.Result {
		houseType := c.parseHousingType(offer)
		if houseType == corporation.HousingTypeUndefined {
			continue
		}

		// enrich offer concurrently
		wg.Add(1)
		go func(offer offerList, houseType corporation.HousingType) {
			defer wg.Done()

			offerDetails, err := c.getOfferDetails(offer.ID)
			if err != nil {
				c.logger.Warn("failed enriching", zap.Any("offer", offer), zap.Error(err), logConnector)
				return
			}

			// add offer to channel
			ch <- c.Map(offerDetails, houseType)
		}(offer, houseType)
	}

	// wait for all offers
	wg.Wait()
	return nil
}

// getOfferDetails info about specific offer housing
func (c *client) getOfferDetails(offerID string) (*offerDetails, error) {
	resp, err := c.Send(offerDetailRequest(offerID))
	if err != nil {
		return nil, err
	}

	var offerDetails struct {
		Result offerDetails `json:"result"`
	}
	if err := json.Unmarshal(resp, &offerDetails); err != nil {
		return nil, fmt.Errorf("error parsing offer detail result %v: %w", string(resp), err)
	}

	return &offerDetails.Result, nil
}

func (c *client) Map(offer *offerDetails, houseType corporation.HousingType) corporation.Offer {
	address := fmt.Sprintf("%s %s-%s %s %s", offer.Street, offer.Housenumber, offer.Housenumberaddition, offer.Postalcode, offer.City.Name)

	numberBedroom, err := strconv.Atoi(offer.Sleepingroom.Amountofrooms)
	if err != nil {
		c.logger.Info("failed parsing number bedroom", zap.String("address", address), zap.Error(err), logConnector)
	}

	// it seems that some appartment from roomspot does not contains rooms while they should (by definition)
	if strings.Contains(offer.Dwellingtype.Localizedname, "studio") {
		numberBedroom = 0
	} else if numberBedroom == 0 {
		numberBedroom = 1
	}

	house := corporation.Housing{
		Type:          houseType,
		Address:       address,
		CityName:      offer.City.Name,
		NumberBedroom: numberBedroom,
		Size:          offer.Areadwelling,
		Price:         offer.Totalrent,
		Garden:        offer.Garden,
		Garage:        true,
		Elevator:      strings.Contains(offer.Dwellingtype.Localizedname, "lift"),
		Balcony:       offer.Balcony,
		Accessible:    true,
	}

	// get address city district
	house.CityDistrict, err = c.mapboxClient.CityDistrictFromAddress(house.Address)
	if err != nil {
		c.logger.Info("could not get city district", zap.String("address", house.Address), zap.Error(err), logConnector)
	}

	// get picture url
	rawPictureURL, err := c.parsePictureURL(offer)
	if err != nil {
		c.logger.Info("failed parsing picture url", zap.Error(err), logConnector)
	}

	return corporation.Offer{
		CorporationName: c.corporation.Name,
		ExternalID:      c.getExternalID(offer),
		Housing:         house,
		URL:             fmt.Sprintf("%s/aanbod/te-huur/details/%s", c.corporation.URL, offer.Urlkey),
		RawPictureURL:   rawPictureURL,
		MinFamilySize:   offer.Minimumhouseholdsize,
		MaxFamilySize:   offer.Maximumhouseholdsize,
		MinAge:          offer.Minimumage,
		MaxAge:          offer.Maximumage,
		MinimumIncome:   offer.Minimumincome,
		MaximumIncome:   offer.Maximumincome,
	}
}

func (c *client) parseHousingType(offer offerList) corporation.HousingType {
	name := strings.ToLower(offer.Dwellingtype.Localizedname)

	if offer.Dwellingtype.Categorie != "woning" || strings.Contains(name, "kamer") {
		return corporation.HousingTypeUndefined
	}

	if strings.Contains(name, "appartement") ||
		strings.Contains(name, "studio") ||
		strings.Contains(name, "flat") {
		return corporation.HousingTypeAppartement
	}

	return corporation.HousingTypeHouse
}

func (c *client) getExternalID(offer *offerDetails) string {
	return fmt.Sprint(offer.Assignmentid) + externalIDSeperator + offer.ID
}

func (c *client) parsePictureURL(offer *offerDetails) (*url.URL, error) {
	if len(offer.Pictures) == 0 {
		return nil, nil
	}

	path := c.corporation.URL + offer.Pictures[0].URI
	pictureURL, err := url.Parse(path)
	if err != nil {
		return nil, fmt.Errorf("failed to parse picture url %s: %w", path, err)
	}

	return pictureURL, nil
}
