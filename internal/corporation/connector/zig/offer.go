package zig

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"

	"github.com/woningfinder/woningfinder/internal/corporation"
	"github.com/woningfinder/woningfinder/pkg/networking"
)

const externalIDSeperator = ";"

type offerList struct {
	Dwellingtype struct {
		Categorie     string `json:"categorie"`
		Localizedname string `json:"localizedName"`
	} `json:"dwellingType"`
	TotalRent     float64 `json:"totalRent"`
	Rentbuy       string  `json:"rentBuy"`
	Iszelfstandig bool    `json:"isZelfstandig"`
	ID            string  `json:"id"`
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
		Categorie           string `json:"categorie"`
		Huurprijsduuractief bool   `json:"huurprijsDuurActief"`
		Localizedname       string `json:"localizedName"`
	} `json:"dwellingType"`
	Totalrent    float64 `json:"totalRent"`
	Sleepingroom struct {
		Amountofrooms string `json:"amountOfRooms"`
		ID            string `json:"id"`
		Localizedname string `json:"localizedName"`
	} `json:"sleepingRoom"`
	Garden                   bool `json:"garden"`
	Balcony                  bool `json:"balcony"`
	Minimumincome            int  `json:"minimumIncome"`
	Maximumincome            int  `json:"maximumIncome"`
	Minimumhouseholdsize     int  `json:"minimumHouseholdSize"`
	Maximumhouseholdsize     int  `json:"maximumHouseholdSize"`
	Minimumage               int  `json:"minimumAge"`
	Maximumage               int  `json:"maximumAge"`
	Inwonendekinderenminimum int  `json:"inwonendeKinderenMinimum"`
	Inwonendekinderenmaximum int  `json:"inwonendeKinderenMaximum"`
	Model                    struct {
		Modelcategorie struct {
			Icon          string `json:"icon"`
			Code          string `json:"code"`
			Toonopwebsite bool   `json:"toonOpWebsite"`
		} `json:"modelCategorie"`
		Incode                            string `json:"inCode"`
		Isvoorextraaanbod                 bool   `json:"isVoorExtraAanbod"`
		Ishospiteren                      bool   `json:"isHospiteren"`
		Advertentiesluitennaeerstereactie bool   `json:"advertentieSluitenNaEersteReactie"`
		Einddatumtonen                    bool   `json:"einddatumTonen"`
		Aantalreactiestonen               bool   `json:"aantalReactiesTonen"`
		Slaagkanstonen                    bool   `json:"slaagkansTonen"`
		ID                                string `json:"id"`
		Localizedname                     string `json:"localizedName"`
	} `json:"model"`
	Rentbuy           string        `json:"rentBuy"`
	Publicationdate   string        `json:"publicationDate"`
	Closingdate       string        `json:"closingDate"`
	Numberofreactions int           `json:"numberOfReactions"`
	Assignmentid      int           `json:"assignmentID"`
	Latitude          string        `json:"latitude"`
	Longitude         string        `json:"longitude"`
	Floorplans        []interface{} `json:"floorplans"`
	Pictures          []struct {
		Label string `json:"label"`
		URI   string `json:"uri"`
		Type  string `json:"type"`
	} `json:"pictures"`
	Areadwelling      float64 `json:"areaDwelling"`
	Areaperceel       string  `json:"areaPerceel"`
	Iszelfstandig     bool    `json:"isZelfstandig"`
	Urlkey            string  `json:"urlKey"`
	Availablefromdate string  `json:"availableFromDate"`
	ID                string  `json:"id"`
	Isgepubliceerd    bool    `json:"isGepubliceerd"`
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

func (c *client) GetOffers() ([]corporation.Offer, error) {
	resp, err := c.Send(offerRequest())
	if err != nil {
		return nil, err
	}

	var result struct {
		Result []offerList `json:"result"`
	}
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("error parsing offer result %v: %w", string(resp), err)
	}

	var offers []corporation.Offer
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
				// do not append the house but logs error
				c.logger.Sugar().Warnf("zig connector: failed enriching %v: %w", offer, err)
				return
			}

			offers = append(offers, c.Map(offerDetails, houseType))
		}(offer, houseType)
	}

	// wait for all offers
	wg.Wait()
	return offers, nil
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
	numberBedroom, err := strconv.Atoi(offer.Sleepingroom.Amountofrooms)
	if err != nil {
		c.logger.Sugar().Infof("zig connector: failed parsing number bedroom: %w", err)
	}

	// it seems that some appartment from roomspot does not contains rooms while they should (by definition)
	if strings.Contains(offer.Dwellingtype.Localizedname, "studio") {
		numberBedroom = 0
	} else if numberBedroom == 0 {
		numberBedroom = 1 // TODO TO FIX
	}

	house := corporation.Housing{
		Type:          houseType,
		Address:       fmt.Sprintf("%s %s-%s %s %s", offer.Street, offer.Housenumber, offer.Housenumberaddition, offer.Postalcode, offer.City.Name),
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
		c.logger.Sugar().Infof("zig connector: could not get city district of %s: %w", house.Address, err)
	}

	// get picture url
	rawPictureURL, err := c.parsePictureURL(offer)
	if err != nil {
		c.logger.Sugar().Info(err)
	}

	return corporation.Offer{
		ExternalID:      c.getExternalID(offer),
		Housing:         house,
		URL:             fmt.Sprintf("%s/aanbod/te-huur/details/%s", c.corporation.URL, offer.Urlkey),
		RawPictureURL:   rawPictureURL,
		SelectionMethod: c.parseSelectionMethod(offer),
		MinFamilySize:   offer.Minimumhouseholdsize,
		MaxFamilySize:   offer.Maximumhouseholdsize,
		MinAge:          offer.Minimumage,
		MaxAge:          offer.Maximumage,
		MinimumIncome:   offer.Minimumincome,
		MaximumIncome:   offer.Maximumincome,
	}
}

func (c *client) parseHousingType(offer offerList) corporation.HousingType {
	if offer.Dwellingtype.Categorie != "woning" || !offer.Iszelfstandig {
		return corporation.HousingTypeUndefined
	}

	name := strings.ToLower(offer.Dwellingtype.Localizedname)
	if strings.Contains(name, "appartement") ||
		strings.Contains(name, "studio") {
		return corporation.HousingTypeAppartement
	}

	return corporation.HousingTypeHouse
}

func (c *client) parseSelectionMethod(offer *offerDetails) corporation.SelectionMethod {
	if offer.Model.Modelcategorie.Code == "inschrijfduur" {
		return corporation.SelectionRegistrationDate
	}

	if offer.Model.Modelcategorie.Code == "reactiedatum" {
		return corporation.SelectionRandom
	}

	return corporation.SelectionFirstComeFirstServed
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
		return nil, fmt.Errorf("zig connector: failed to parse picture url %s: %w", path, err)
	}

	return pictureURL, nil
}
