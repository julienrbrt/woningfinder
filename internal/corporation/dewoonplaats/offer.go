package dewoonplaats

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/woningfinder/woningfinder/internal/domain/entity"
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
	} `json:"woningen"`
}

func (c *client) FetchOffer() ([]entity.Offer, error) {
	req, err := offerRequest()
	if err != nil {
		return nil, err
	}

	resp, err := c.Send(req)
	if err != nil {
		return nil, err
	}

	var result offerResult
	if err := json.Unmarshal(resp.Result, &result); err != nil {
		return nil, fmt.Errorf("error parsing login result %v: %w", string(resp.Result), err)
	}

	var offers []entity.Offer
	for _, house := range result.Offer {
		houseType := c.parseHousingType(house.HousingType)

		if !house.ForRent || houseType == entity.HousingTypeUndefined || house.RentPrice == 0 {
			continue
		}

		cityDistrict := house.District
		if cityDistrict == "" {
			cityDistrict, err = c.mapboxClient.CityDistrictFromCoords(fmt.Sprintf("%f", house.Latitude), fmt.Sprintf("%f", house.Longitude))
			if err != nil {
				c.logger.Sugar().Infof("could not get city district of %s: %w", house.Address, err)
			}
		}

		newHouse := entity.Housing{
			Type: entity.HousingType{
				Type: houseType,
			},
			Address: fmt.Sprintf("%s %s %s", house.Address, house.Postcode, house.City),
			City: entity.City{
				Name: house.City,
			},
			CityDistrict: entity.CityDistrict{
				CityName: house.City,
				Name:     cityDistrict,
			},
			Latitude:         house.Latitude,
			Longitude:        house.Longitude,
			EnergieLabel:     house.EnergieLabel,
			NumberRoom:       len(house.RoomSize),
			NumberBedroom:    house.NumberBedroom,
			Size:             c.parseHouseSize(house.Size),
			Price:            house.RentPrice,
			BuildingYear:     house.BuildingYear,
			HousingAllowance: house.HasLowRentPrice && !house.RentLuxe && len(house.RentPriceForAllowance) > 0,
			Garden:           len(house.Garden) > 0,
			Garage:           house.Garage,
			Elevator:         house.Lift,
			Balcony:          house.Balcony,
			Attic:            house.Attic,
			Accessible:       house.Accessible,
		}

		offer := entity.Offer{
			ExternalID: house.ID,
			Housing:    newHouse,
			URL:        fmt.Sprintf("https://www.dewoonplaats.nl/ik-zoek-woonruimte/!/woning/%s/", house.ID),
			SelectionMethod: entity.SelectionMethod{
				Method: c.parseSelectionMethod(house.IsSelectionRandom),
			},
			SelectionDate: c.parseSelectionDate(house.SelectionDate),
			MinIncome:     house.Criteria.MinInkomen,
			MaxIncome:     house.Criteria.MaxInkomen,
			MinFamilySize: house.Criteria.MinGezinsgrootte,
			MaxFamilySize: house.Criteria.MaxGezinsgrootte,
			MinAge:        house.Criteria.MinLeeftijd,
			MaxAge:        house.Criteria.MaxLeeftijd,
		}

		offers = append(offers, offer)
	}

	return offers, nil
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

func (c *client) parseHousingType(houseType []string) entity.Type {
	if len(houseType) == 0 {
		return entity.HousingTypeUndefined
	}

	for _, h := range houseType {
		h = strings.ToLower(h)
		if h == "appartement" {
			return entity.HousingTypeAppartement
		} else if h == "eengezinswoning" {
			return entity.HousingTypeHouse
		}
	}

	return entity.HousingTypeUndefined
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

func (c *client) parseSelectionMethod(random bool) entity.Method {
	if random {
		return entity.SelectionRandom
	}

	return entity.SelectionFirstComeFirstServed
}
