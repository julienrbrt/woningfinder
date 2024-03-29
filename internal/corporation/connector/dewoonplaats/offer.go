package dewoonplaats

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/julienrbrt/woningfinder/internal/corporation"
	"github.com/julienrbrt/woningfinder/pkg/networking"
	"go.uber.org/zap"
)

const methodOffer = "ZoekWoningen"

type offer struct {
	ID          string   `json:"id"`
	HousingType []string `json:"soort"`
	Criteria    struct {
		KinderenValid    bool `json:"kinderen_valid"`
		MaxGezinsgrootte int  `json:"max_gezinsgrootte"`
		MaxInkomen       int  `json:"max_inkomen"`
		MaxLeeftijd      int  `json:"max_leeftijd"`
		MinGezinsgrootte int  `json:"min_gezinsgrootte"`
		MinInkomen       int  `json:"min_inkomen"`
		MinLeeftijd      int  `json:"min_leeftijd"`
	} `json:"criteria"`
	Address         string  `json:"adres"`
	District        string  `json:"wijk"`
	CityName        string  `json:"plaats"`
	Postcode        string  `json:"postcode"`
	RentPrice       float64 `json:"relevante_huurprijs,omitempty"`
	NumberBedroom   int     `json:"slaapkamers"`
	CV              bool    `json:"cv"`
	Balcony         bool    `json:"balkon"`
	Garage          bool    `json:"garage"`
	ForRent         bool    `json:"ishuur"`
	HasLowRentPrice bool    `json:"ishuurlaag"`
	Lift            bool    `json:"lift"`
	Garden          string  `json:"tuin"`
	Size            string  `json:"woonoppervlak"`
	Accessible      bool    `json:"rolstoeltoegankelijk"`
	Thumbnail       string  `json:"thumbnail"`
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

func (c *client) FetchOffers(ch chan<- corporation.Offer) error {
	req, err := offerRequest()
	if err != nil {
		return err
	}

	resp, err := c.Send(req)
	if err != nil {
		return err
	}

	var result struct {
		Offer []offer `json:"woningen"`
	}

	if err := json.Unmarshal(resp.Result, &result); err != nil {
		return fmt.Errorf("error parsing offer result %v: %w", string(resp.Result), err)
	}

	for _, offer := range result.Offer {
		houseType := c.parseHousingType(offer.HousingType)

		if !offer.ForRent || houseType == corporation.HousingTypeUndefined || offer.RentPrice == 0 {
			continue
		}

		// add offer to channel
		ch <- c.Map(offer, houseType)
	}

	return nil
}

func (c *client) Map(offer offer, houseType corporation.HousingType) corporation.Offer {
	house := corporation.Housing{
		Type:          houseType,
		Address:       fmt.Sprintf("%s %s %s", offer.Address, offer.Postcode, offer.CityName),
		CityName:      offer.CityName,
		NumberBedroom: offer.NumberBedroom,
		Size:          c.parseHouseSize(offer.Size),
		Price:         offer.RentPrice,
		Garden:        len(offer.Garden) > 0,
		Garage:        offer.Garage,
		Elevator:      offer.Lift,
		Balcony:       offer.Balcony,
		Accessible:    offer.Accessible,
	}

	// get address city district
	var err error
	house.CityDistrict, err = c.mapboxClient.CityDistrictFromAddress(house.Address)
	if err != nil {
		house.CityDistrict = offer.District
		c.logger.Info("could not get city district", zap.String("address", house.Address), zap.Error(err), logConnector)
	}

	rawPictureURL, err := c.parsePictureURL(offer.Thumbnail)
	if err != nil {
		c.logger.Info("failed parsing picture url", zap.Error(err), logConnector)
	}

	return corporation.Offer{
		CorporationName: c.corporation.Name,
		ExternalID:      offer.ID,
		Housing:         house,
		URL:             fmt.Sprintf("https://www.dewoonplaats.nl/ik-zoek-woonruimte/!/woning/%s/", offer.ID),
		RawPictureURL:   rawPictureURL,
		MinFamilySize:   offer.Criteria.MinGezinsgrootte,
		MaxFamilySize:   offer.Criteria.MaxGezinsgrootte,
		MinAge:          offer.Criteria.MinLeeftijd,
		MaxAge:          offer.Criteria.MaxLeeftijd,
		MinimumIncome:   offer.Criteria.MinInkomen,
		MaximumIncome:   offer.Criteria.MaxInkomen,
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
	size, _ := strconv.ParseFloat(strings.ReplaceAll(houseSize, ",", "."), 32)

	return size
}

func (c *client) parsePictureURL(path string) (*url.URL, error) {
	if path == "" {
		return nil, nil
	}

	pictureURL, err := url.Parse(c.corporation.URL + path)
	if err != nil {
		return nil, fmt.Errorf("failed to parse picture url %s: %w", path, err)
	}

	return pictureURL, nil
}
