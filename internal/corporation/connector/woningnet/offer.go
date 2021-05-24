package woningnet

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
	"github.com/woningfinder/woningfinder/internal/corporation"
	"github.com/woningfinder/woningfinder/internal/corporation/city"
	"github.com/woningfinder/woningfinder/pkg/networking"
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
	Address                     string      `json:"Adres"`
	CityAndDistrict             string      `json:"PlaatsWijk"`
	Omschrijving                string      `json:"Omschrijving"`
	Aanbieder                   string      `json:"Aanbieder"`
	Prijs                       string      `json:"Prijs"`
	Kamers                      string      `json:"Kamers"`
	AfbeeldingURL               string      `json:"AfbeeldingUrl"`
	SoortWoning                 interface{} `json:"SoortWoning"`
	Slaagkans                   interface{} `json:"Slaagkans"`
	SlaagkansText               interface{} `json:"SlaagkansText"`
	SelectionDate               time.Time   `json:"PublicatieEinddatum"`
	PublicatieBegindatum        time.Time   `json:"PublicatieBegindatum"`
	PublicatieEinddatumVolledig string      `json:"PublicatieEinddatumVolledig"`
	PublicatieWachttijd         string      `json:"PublicatieWachttijd"`
	PublicatieBeschikbaarPer    string      `json:"PublicatieBeschikbaarPer"`
	IsToonTegels                bool        `json:"IsToonTegels"`
	Status                      string      `json:"Status"`
	AantalReacties              interface{} `json:"AantalReacties"`
	VoorlopigePositie           interface{} `json:"VoorlopigePositie"`
	AantalWoningen              interface{} `json:"AantalWoningen"`
	AantalBeschikbaar           interface{} `json:"AantalBeschikbaar"`
	PublicatieID                string      `json:"PublicatieId"`
	ResterendetijdUur           float64     `json:"ResterendetijdUur"`
	ResterendetijdDagen         int         `json:"ResterendetijdDagen"`
	ShowResterendeTijd          bool        `json:"ShowResterendeTijd"`
	ResterendeTijdMinderDanUur  bool        `json:"ResterendeTijdMinderDanUur"`
	AdvertentieURL              string      `json:"AdvertentieUrl"`
	MinimalePositie             interface{} `json:"MinimalePositie"`
	ResultaatIndex              int         `json:"ResultaatIndex"`
	PreviewURL                  string      `json:"PreviewUrl"`
	Latitude                    float64     `json:"Latitude"`
	Longitude                   float64     `json:"Longitude"`
	MapsIcon                    struct {
		Origin struct {
			X string `json:"X"`
			Y string `json:"Y"`
		} `json:"Origin"`
		Size struct {
			X string `json:"X"`
			Y string `json:"Y"`
		} `json:"Size"`
		CSSClass string `json:"CssClass"`
		Label    string `json:"Label"`
	} `json:"MapsIcon"`
	ToonMessageList                 bool   `json:"ToonMessageList"`
	PublicatieModel                 string `json:"PublicatieModel"`
	CurrentRegioCode                string `json:"CurrentRegioCode"`
	IconName                        string `json:"IconName"`
	PrijsHelpText                   string `json:"PrijsHelpText"`
	IsNieuw                         bool   `json:"IsNieuw"`
	IsWoonwensMatch                 bool   `json:"IsWoonwensMatch"`
	Woonoppervlakte                 string `json:"Woonoppervlakte"`
	DetailSoortCode                 string `json:"DetailSoortCode"`
	HouseType                       string `json:"DetailSoortOmschrijving"`
	PublicatieModulesCode           string `json:"PublicatieModulesCode"`
	RegioCode                       string `json:"RegioCode"`
	PublicatieDetailModel           string `json:"PublicatieDetailModel"`
	IsKoop                          bool   `json:"IsKoop"`
	IsSocialeHuur                   bool   `json:"IsSocialeHuur"`
	IsVrijeSectorhuur               bool   `json:"IsVrijeSectorhuur"`
	IsVrijesectorOfKoop             bool   `json:"IsVrijesectorOfKoop"`
	IsLoting                        bool   `json:"IsLoting"`
	IsGarage                        bool   `json:"IsGarage"`
	IsTeWoon                        bool   `json:"IsTeWoon"`
	IsDirectTeHuur                  bool   `json:"IsDirectTeHuur"`
	WoningType                      int    `json:"WoningType"`
	WoningTypeCSSClass              string `json:"WoningTypeCssClass"`
	Samenvatting                    string `json:"Samenvatting"`
	Toegankelijkheidslabel          string `json:"Toegankelijkheidslabel"`
	ToegankelijkheidslabelCSSClass  string `json:"ToegankelijkheidslabelCssClass"`
	ToegankelijkheidslabelHelpTekst string `json:"ToegankelijkheidslabelHelpTekst"`
	HasFlexibelhurenIndicator       bool   `json:"HasFlexibelhurenIndicator"`
	FlexibelhurenIndicatorCSSClass  string `json:"FlexibelhurenIndicatorCssClass"`
	ToonSlaagkans                   bool   `json:"ToonSlaagkans"`
}

func offerRequest(commandURL, command string) (networking.Request, error) {
	req := request{
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

func (c *client) GetOffers() ([]corporation.Offer, error) {
	var offers []corporation.Offer

	for _, selectionMethod := range c.corporation.SelectionMethod {
		req, err := offerRequest(selectionMethodCommandMap[selectionMethod], "")
		if err != nil {
			return nil, err
		}

		resp, err := c.Send(req)
		if err != nil {
			return nil, err
		}

		var response response
		if err := json.Unmarshal(resp, &response); err != nil {
			return nil, fmt.Errorf("error parsing offer response %v: %w", string(resp), err)
		}

		paginatedOffers, err := c.getPaginatedOffers(selectionMethodCommandMap[selectionMethod], response)
		if err != nil {
			c.logger.Sugar().Warnf("woningnet connector: error getting paginated offers: %w", err)
		}
		respOffers := append(response.Offer, paginatedOffers...)

		for _, offer := range respOffers {
			houseType := c.parseHousingType(offer.HouseType)
			if houseType == corporation.HousingTypeUndefined {
				continue
			}

			result, err := c.Map(offer, houseType, selectionMethod)
			if err != nil {
				c.logger.Sugar().Warnf("woningnet connector: error getting parsing %s: %w", offer.Address, err)
				continue
			}

			offers = append(offers, result)
		}
	}

	return offers, nil
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

func (c *client) Map(offer offer, houseType corporation.HousingType, selectionMethod corporation.SelectionMethod) (corporation.Offer, error) {
	var err error
	var minFamilySize, maxFamilySize, minAge, maxAge int
	var offerURL = c.corporation.URL + offer.AdvertentieURL
	var cityName = strings.Split(offer.CityAndDistrict, " - ")

	// fill in already known housing characteristics
	house := corporation.Housing{
		Type:    houseType,
		Address: fmt.Sprintf("%s %s", offer.Address, cityName[0]),
		City: corporation.City{
			Name: cityName[0],
		},
		Accessible: offer.ToegankelijkheidslabelCSSClass != "ToegankelijkheidslabelZON",
	}

	house.Price, err = c.parsePrice(offer.Prijs)
	if err != nil {
		return corporation.Offer{}, fmt.Errorf("error mapping offer %s: %w", offerURL, err)
	}

	house.Size, err = c.parseHouseSize(offer.Woonoppervlakte)
	if err != nil {
		return corporation.Offer{}, fmt.Errorf("error mapping offer %s: %w", offerURL, err)
	}

	// get address city district
	if city.HasSuggestedCityDistrict(house.City.Name) {
		house.CityDistrict, err = c.mapboxClient.CityDistrictFromAddress(house.Address)
		if err != nil {
			house.CityDistrict = cityName[len(cityName)-1]
			c.logger.Sugar().Infof("woningnet connector: could not get city district of %s: %w", house.Address, err)
		}
	}

	// TODO
	// c.collector.OnHTML("#Overzicht", func(e *colly.HTMLElement) {
	// 	// minFamilySize, maxFamilySize, minAge, maxAge int
	// 	// fmt.Println(e)
	// })

	c.collector.OnHTML("#Kenmerken", func(e *colly.HTMLElement) {
		table := e.DOM.ChildrenFiltered(".contentBlock")

		// energy label
		house.EnergyLabel = c.getContentValue("Energielabel", table)

		// building year
		house.BuildingYear, err = strconv.Atoi(c.getContentValue("Bouwjaar", table))
		if err != nil {
			c.logger.Sugar().Infof("woningnet connector: error parsing building year of %s: %w", house.Address, err)
		}

		// number bedroom
		house.NumberBedroom, err = c.parseBedroom(c.getContentValue("Aantal kamers (incl. woonkamer)", table))
		if err != nil {
			c.logger.Sugar().Infof("woningnet connector: error parsing number bedroom of %s: %w", house.Address, err)
		}

		// attic
		house.Attic = len(c.getContentValue("Zolder", table)) > 0

		// outside
		outside := c.getContentValue("Buitenruimte", table)
		house.Balcony = strings.Contains(outside, "balkon") || strings.Contains(outside, "terras")
		house.Garden = strings.Contains(outside, "tuin")
		house.Garage = strings.Contains(outside, "garage") || strings.Contains(outside, "parkeer")

		// TODO add lift parsing

	})

	if err = c.collector.Visit(offerURL); err != nil {
		return corporation.Offer{}, fmt.Errorf("error visiting %s: %w", offerURL, err)
	}

	return corporation.Offer{
		ExternalID:      fmt.Sprintf("%s/%s", offer.PublicatieID, offer.CurrentRegioCode),
		Housing:         house,
		URL:             offerURL,
		SelectionMethod: selectionMethod,
		SelectionDate:   offer.SelectionDate,
		MinFamilySize:   minFamilySize,
		MaxFamilySize:   maxFamilySize,
		MinAge:          minAge,
		MaxAge:          maxAge,
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
	price, err := strconv.ParseFloat(strings.ReplaceAll(strings.TrimLeft(priceStr, "€ "), ",", "."), 64)
	if err != nil {
		return 0, fmt.Errorf("error parsing housing price %s: %w", priceStr, err)
	}

	return price, nil
}

func (c *client) parseHouseSize(sizeStr string) (float64, error) {
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

// parseContentValue is use to get the value from a table in the woningnet housing description page
func (c *client) getContentValue(name string, selection *goquery.Selection) string {
	var value string
	selection.Children().Each(func(i int, s *goquery.Selection) {
		if strings.EqualFold(s.Children().First().Text(), name) {
			value = s.Children().Next().Text()
		}
	})

	return strings.ToLower(value)
}
