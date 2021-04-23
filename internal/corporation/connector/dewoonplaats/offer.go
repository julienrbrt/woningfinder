package dewoonplaats

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/woningfinder/woningfinder/internal/corporation"
	"github.com/woningfinder/woningfinder/pkg/networking"
)

const methodOffer = "ZoekWoningen"

type offer struct {
	ID          string   `json:"id"`
	HousingType []string `json:"soort"`
	Criteria    struct {
		KinderenValid    bool   `json:"kinderen_valid"`
		MaxGezinsgrootte int    `json:"max_gezinsgrootte"`
		MaxInkomen       int    `json:"max_inkomen"`
		MaxLeeftijd      int    `json:"max_leeftijd"`
		MinGezinsgrootte int    `json:"min_gezinsgrootte"`
		MinInkomen       int    `json:"min_inkomen"`
		MinLeeftijd      int    `json:"min_leeftijd"`
		Omschrijving     string `json:"omschrijving"`
	} `json:"criteria"`
	Latitude              float64 `json:"lat"`
	Longitude             float64 `json:"lng"`
	Address               string  `json:"adres"`
	District              string  `json:"wijk"`
	City                  string  `json:"plaats"`
	Postcode              string  `json:"postcode"`
	RentPrice             float64 `json:"relevante_huurprijs,omitempty"`
	RentPriceForAllowance string  `json:"toeslagprijs"`
	RentLuxe              bool    `json:"tehuur_luxehuur,omitempty"`
	MapsURL               string  `json:"mapslink"`
	BuildingYear          int     `json:"bouwjaar"`
	EnergieLabel          string  `json:"energielabel"`
	NumberBedroom         int     `json:"slaapkamers"`
	CV                    bool    `json:"cv"`
	Balcony               bool    `json:"balkon"`
	Garage                bool    `json:"garage"`
	Historic              bool    `json:"historic"`
	ForRent               bool    `json:"ishuur"`
	HasLowRentPrice       bool    `json:"ishuurlaag"`
	Lift                  bool    `json:"lift"`
	Garden                string  `json:"tuin"`
	Attic                 bool    `json:"zolder"`
	SelectionDate         string  `json:"lotingsdatum"`
	IsSelectionRandom     bool    `json:"loting"`
	Size                  string  `json:"woonoppervlak"`
	Accessible            bool    `json:"rolstoeltoegankelijk"`
	RoomSize              []struct {
		Name string `json:"titel"`
		Size string `json:"oppervlak"`
	} `json:"vertrekken"`
}

func offerRequest() (networking.Request, error) {
	req := request{
		ID:     1,
		Method: methodOffer,
		Params: []interface{}{
			struct {
				ForRent bool `json:"tehuur"`
			}{
				ForRent: true,
			},
			"",
			true,
		},
	}

	body, err := json.Marshal(req)
	if err != nil {
		return networking.Request{}, fmt.Errorf("error while marshaling %v: %w", req, err)
	}

	request := networking.Request{
		Path:   "/woonplaats_digitaal/woonvinder",
		Method: http.MethodPost,
		Body:   bytes.NewBuffer(body),
	}

	return request, nil
}

func (c *client) GetOffers() ([]corporation.Offer, error) {
	req, err := offerRequest()
	if err != nil {
		return nil, err
	}

	resp, err := c.Send(req)
	if err != nil {
		return nil, err
	}

	var result struct {
		Offer []offer `json:"woningen"`
	}

	if err := json.Unmarshal(resp.Result, &result); err != nil {
		return nil, fmt.Errorf("error parsing login result %v: %w", string(resp.Result), err)
	}

	var offers []corporation.Offer
	for _, offer := range result.Offer {
		houseType := c.parseHousingType(offer.HousingType)

		if !offer.ForRent || houseType == corporation.HousingTypeUndefined || offer.RentPrice == 0 {
			continue
		}

		offers = append(offers, c.Map(offer, houseType))
	}

	return offers, nil
}

func (c *client) Map(offer offer, houseType corporation.HousingType) corporation.Offer {
	house := corporation.Housing{
		Type:    houseType,
		Address: fmt.Sprintf("%s %s %s", offer.Address, offer.Postcode, offer.City),
		City: corporation.City{
			Name: offer.City,
		},
		EnergieLabel:  offer.EnergieLabel,
		NumberRoom:    len(offer.RoomSize),
		NumberBedroom: offer.NumberBedroom,
		Size:          c.parseHouseSize(offer.Size),
		Price:         offer.RentPrice,
		BuildingYear:  offer.BuildingYear,
		Garden:        len(offer.Garden) > 0,
		Garage:        offer.Garage,
		Elevator:      offer.Lift,
		Balcony:       offer.Balcony,
		Attic:         offer.Attic,
		Accessible:    offer.Accessible,
	}

	// get address city district
	var err error
	house.CityDistrict, err = c.mapboxClient.CityDistrictFromAddress(house.Address)
	if err != nil {
		house.CityDistrict = offer.District
		c.logger.Sugar().Infof("could not get city district of %s: %w", house.Address, err)
	}

	return corporation.Offer{
		ExternalID:      offer.ID,
		Housing:         house,
		URL:             fmt.Sprintf("https://www.dewoonplaats.nl/ik-zoek-woonruimte/!/woning/%s/", offer.ID),
		SelectionMethod: c.parseSelectionMethod(offer.IsSelectionRandom),
		SelectionDate:   c.parseSelectionDate(offer.SelectionDate),
		MinFamilySize:   offer.Criteria.MinGezinsgrootte,
		MaxFamilySize:   offer.Criteria.MaxGezinsgrootte,
		MinAge:          offer.Criteria.MinLeeftijd,
		MaxAge:          offer.Criteria.MaxLeeftijd,
	}
}

func (c *client) parseHousingType(houseType []string) corporation.HousingType {
	if len(houseType) == 0 {
		return corporation.HousingTypeUndefined
	}

	for _, h := range houseType {
		h = strings.ToLower(h)
		if h == "appartement" {
			return corporation.HousingTypeAppartement
		} else if h == "eengezinswoning" {
			return corporation.HousingTypeHouse
		}
	}

	return corporation.HousingTypeUndefined
}

func (c *client) parseHouseSize(houseSize string) float64 {
	size, _ := strconv.ParseFloat(strings.ReplaceAll(houseSize, ",", "."), 16)

	return size
}

func (c *client) parseSelectionDate(str string) time.Time {
	if len(str) == 0 {
		return time.Time{}
	}

	date, err := time.Parse(time.RFC3339, str)
	if err != nil {
		return time.Time{}
	}

	return date
}

func (c *client) parseSelectionMethod(random bool) corporation.SelectionMethod {
	if random {
		return corporation.SelectionRandom
	}

	return corporation.SelectionFirstComeFirstServed
}
