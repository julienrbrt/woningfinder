package dewoonplaats

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/woningfinder/woningfinder/internal/corporation"
	"github.com/woningfinder/woningfinder/pkg/networking"
)

const methodOffer = "ZoekWoningen"

type offerResult struct {
	Total int `json:"aantal"`
	Offer []struct {
		Address      string `json:"adres"`
		Balcony      bool   `json:"balkon"`
		BuildingYear int    `json:"bouwjaar"`
		Criteria     struct {
			KinderenValid    bool   `json:"kinderen_valid"`
			MaxGezinsgrootte int    `json:"max_gezinsgrootte"`
			MaxInkomen       int    `json:"max_inkomen"`
			MaxLeeftijd      int    `json:"max_leeftijd"`
			MinGezinsgrootte int    `json:"min_gezinsgrootte"`
			MinInkomen       int    `json:"min_inkomen"`
			MinLeeftijd      int    `json:"min_leeftijd"`
			Omschrijving     string `json:"omschrijving"`
		} `json:"criteria"`
		Cv                      bool     `json:"cv"`
		Datum                   string   `json:"datum"`
		EnergieLabel            string   `json:"energielabel"`
		Etage                   string   `json:"etage"`
		Foto                    string   `json:"foto"`
		Garage                  bool     `json:"garage"`
		Historic                bool     `json:"historic"`
		ID                      string   `json:"id"`
		ForRent                 bool     `json:"ishuur"`
		HasLowRentPrice         bool     `json:"ishuurlaag"`
		Keuken                  string   `json:"keuken"`
		Latitude                float64  `json:"lat"`
		Lift                    bool     `json:"lift"`
		Longitude               float64  `json:"lng"`
		Loting                  bool     `json:"loting"`
		Lotingsdatum            string   `json:"lotingsdatum"`
		CanApply                bool     `json:"magreageren"`
		HasAppliedOn            string   `json:"gereageerd_op"`
		MapsURL                 string   `json:"mapslink"`
		City                    string   `json:"plaats"`
		Postcode                string   `json:"postcode"`
		Reactiedatum            string   `json:"reactiedatum"`
		Reacties                int      `json:"reacties"`
		Recreatieruimte         bool     `json:"recreatieruimte"`
		RentPrice               float64  `json:"relevante_huurprijs,omitempty"`
		AccessibilityScooter    bool     `json:"rollatortoegankelijk"`
		AccessibilityWheelchair bool     `json:"rolstoeltoegankelijk"`
		Servicekosten           string   `json:"servicekosten"`
		NumberBedroom           int      `json:"slaapkamers"`
		HousingType             []string `json:"soort"`
		Straat                  string   `json:"straat"`
		TehuurLuxehuur          bool     `json:"tehuur_luxehuur,omitempty"`
		Thumbnail               string   `json:"thumbnail"`
		Toeslagprijs            string   `json:"toeslagprijs"`
		Garden                  string   `json:"tuin"`
		Verbruikskosten         string   `json:"verbruikskosten"`
		Vertrekken              []struct {
			Oppervlak string `json:"oppervlak"`
			Titel     string `json:"titel"`
		} `json:"vertrekken"`
		Verwarming      string `json:"verwarming"`
		Wijk            string `json:"wijk"`
		Wijkid          int    `json:"wijkid"`
		Wijzigingsdatum string `json:"wijzigingsdatum"`
		SizeM2          string `json:"woonoppervlak"`
		WrdID           int    `json:"wrd_id,omitempty"`
		Attic           bool   `json:"zolder"`
	} `json:"woningen"`
}

func (c *client) FetchOffer(minimumPrice int) ([]corporation.Offer, error) {
	req, err := c.offerRequest(minimumPrice)
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
		houseType := c.getHousingType(house.HousingType)

		if !house.ForRent || houseType == corporation.Undefined || house.RentPrice == 0 {
			continue
		}

		houseSizeM2, err := strconv.Atoi(house.SizeM2)
		if err != nil {
			log.Printf(fmt.Errorf("error while parsing size of house %v: %w", house, err).Error())
			houseSizeM2 = 0
		}

		newHouse := corporation.Housing{
			Type:    houseType,
			Address: fmt.Sprintf("%s %s %s", house.Address, house.Postcode, house.City),
			Location: corporation.Location{
				Latitude:  house.Latitude,
				Longitude: house.Longitude,
			},
			EnergieLabel:            house.EnergieLabel,
			NumberBedroom:           house.NumberBedroom,
			SizeM2:                  houseSizeM2,
			Price:                   house.RentPrice,
			BuildingYear:            house.BuildingYear,
			HousingAllowance:        house.HasLowRentPrice && len(house.Toeslagprijs) > 0,
			Garden:                  len(house.Garden) > 0,
			Garage:                  house.Garage,
			Elevator:                house.Lift,
			Balcony:                 house.Balcony,
			AccessibilityScooter:    house.AccessibilityScooter,
			AccessibilityWheelchair: house.AccessibilityWheelchair,
			Attic:                   house.Attic,
			Historic:                house.Historic,
		}

		offer := corporation.Offer{
			URL:     &url.URL{Scheme: "https", Host: "www.dewoonplaats.nl", Path: fmt.Sprintf("ik-zoek-woonruimte/!/woning/%s/", house.ID)},
			Housing: newHouse,
			City: corporation.City{
				Name: house.City,
			},
		}

		offers = append(offers, offer)
	}

	return offers, nil
}

func (c *client) offerRequest(minimumPrice int) (networking.Request, error) {
	req := request{
		ID:     1,
		Method: methodOffer,
		Params: []interface{}{
			struct {
				MinimumPrice int  `json:"prijsvanaf"`
				ForRent      bool `json:"tehuur"`
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

func (c *client) getHousingType(houseType []string) corporation.HousingType {
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
