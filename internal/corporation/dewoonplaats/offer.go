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

type offerResult struct {
	Offer []struct {
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
		Latitude                float64 `json:"lat"`
		Longitude               float64 `json:"lng"`
		Address                 string  `json:"adres"`
		District                string  `json:"wijk"`
		City                    string  `json:"plaats"`
		Postcode                string  `json:"postcode"`
		RentPrice               float64 `json:"relevante_huurprijs,omitempty"`
		RentPriceForAllowance   string  `json:"toeslagprijs"`
		RentLuxe                bool    `json:"tehuur_luxehuur,omitempty"`
		MapsURL                 string  `json:"mapslink"`
		BuildingYear            int     `json:"bouwjaar"`
		EnergieLabel            string  `json:"energielabel"`
		NumberBedroom           int     `json:"slaapkamers"`
		CV                      bool    `json:"cv"`
		Balcony                 bool    `json:"balkon"`
		Garage                  bool    `json:"garage"`
		Historic                bool    `json:"historic"`
		ForRent                 bool    `json:"ishuur"`
		HasLowRentPrice         bool    `json:"ishuurlaag"`
		Lift                    bool    `json:"lift"`
		AccessibilityScooter    bool    `json:"rollatortoegankelijk"`
		AccessibilityWheelchair bool    `json:"rolstoeltoegankelijk"`
		Garden                  string  `json:"tuin"`
		Attic                   bool    `json:"zolder"`
		CanApply                bool    `json:"magreageren"`
		HasAppliedOn            string  `json:"gereageerd_op"`
		SelectionDate           string  `json:"lotingsdatum"`
		IsSelectionRandom       bool    `json:"loting"`
		Size                    string  `json:"woonoppervlak"`
		RoomSize                []struct {
			Name string `json:"titel"`
			Size string `json:"oppervlak"`
		} `json:"vertrekken"`
	} `json:"woningen"`
}

func (c *client) FetchOffer(minimumPrice float64) ([]corporation.Offer, error) {
	req, err := offerRequest(minimumPrice)
	if err != nil {
		return nil, err
	}

	resp, err := c.Send(req)
	if err != nil {
		return nil, err
	}

	var result offerResult
	if err := json.Unmarshal(resp.Result, &result); err != nil {
		return nil, fmt.Errorf("error parsing login result %v: %w", resp.Result, err)
	}

	var offers []corporation.Offer
	for _, house := range result.Offer {
		houseType := c.parseHousingType(house.HousingType)

		if !house.ForRent || houseType == corporation.Undefined || house.RentPrice == 0 {
			continue
		}

		newHouse := corporation.Housing{
			Type: corporation.HousingType{
				Type: houseType,
			},
			Address: fmt.Sprintf("%s %s %s", house.Address, house.Postcode, house.City),
			City: corporation.City{
				Name: house.City,
				District: []corporation.District{
					{Name: house.District},
				},
			},
			Latitude:                house.Latitude,
			Longitude:               house.Longitude,
			EnergieLabel:            house.EnergieLabel,
			NumberRoom:              len(house.RoomSize),
			NumberBedroom:           house.NumberBedroom,
			Size:                    c.parseHouseSize(houseType, house.Size),
			Price:                   house.RentPrice,
			BuildingYear:            house.BuildingYear,
			HousingAllowance:        house.HasLowRentPrice && !house.RentLuxe && len(house.RentPriceForAllowance) > 0,
			Garden:                  len(house.Garden) > 0,
			CV:                      house.CV,
			Garage:                  house.Garage,
			Elevator:                house.Lift,
			Balcony:                 house.Balcony,
			AccessibilityScooter:    house.AccessibilityScooter,
			AccessibilityWheelchair: house.AccessibilityWheelchair,
			Attic:                   house.Attic,
			Historic:                house.Historic,
		}

		offer := corporation.Offer{
			ExternalID: house.ID,
			Housing:    newHouse,
			URL:        fmt.Sprintf("https://www.dewoonplaats.nl/ik-zoek-woonruimte/!/woning/%s/", house.ID),
			SelectionMethod: corporation.SelectionMethod{
				Method: c.parseSelectionMethod(house.IsSelectionRandom),
			},
			SelectionDate: c.parseSelectionDate(house.SelectionDate),
			CanApply:      house.CanApply,
			HasApplied:    len(house.HasAppliedOn) > 0,
		}

		offers = append(offers, offer)
	}

	return offers, nil
}

func offerRequest(minimumPrice float64) (networking.Request, error) {
	req := request{
		ID:     1,
		Method: methodOffer,
		Params: []interface{}{
			struct {
				MinimumPrice float64 `json:"prijsvanaf"`
				ForRent      bool    `json:"tehuur"`
			}{
				MinimumPrice: minimumPrice,
				ForRent:      true,
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

func (c *client) parseHousingType(houseType []string) corporation.Type {
	if len(houseType) == 0 {
		return corporation.Undefined
	}

	for _, h := range houseType {
		h = strings.ToLower(h)
		if h == "parkeren" {
			return corporation.Parking
		} else if h == "appartement" {
			return corporation.Appartement
		} else if h == "eengezinswoning" {
			return corporation.House
		}
	}

	return corporation.Undefined
}

func (c *client) parseHouseSize(houseType corporation.Type, houseSize string) float64 {
	var size float64
	if houseType == corporation.House || houseType == corporation.Appartement {
		size, _ = strconv.ParseFloat(strings.ReplaceAll(houseSize, ",", "."), 32)
	}

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

func (c *client) parseSelectionMethod(random bool) corporation.Method {
	if random {
		return corporation.SelectionRandom
	}

	return corporation.SelectionFirstComeFirstServed
}
