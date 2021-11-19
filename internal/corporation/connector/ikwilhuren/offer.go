package ikwilhuren

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/woningfinder/woningfinder/internal/corporation"
	"go.uber.org/zap"
)

const offerReserved = "Onder optie"

func (c *client) FetchOffers(ch chan<- corporation.Offer) error {
	defer close(ch)

	offers := map[string]*corporation.Offer{}

	// create another collector for housing details
	paginationCollector := c.collector.Clone()
	detailCollector := c.collector.Clone()

	// check paginating
	c.collector.OnHTML("ul.pagination", func(el *colly.HTMLElement) {
		pageMax, err := strconv.Atoi(el.ChildText("li.page-item > a.page-link > span.sr-only"))
		if err != nil {
			c.logger.Warn("error while parsing pagination", zap.Error(err), logConnector)
		}

		// starts from page 2
		for i := 2; i < pageMax; i++ {
			// visit other pages
			paginatedURL := fmt.Sprintf("%s/pagina/%d?action=epl_search&post_type=rental", c.corporation.URL, i)
			if err := paginationCollector.Visit(paginatedURL); err != nil {
				c.logger.Info("error while checking pagination", zap.String("url", paginatedURL), zap.Error(err), logConnector)
			}
		}
	})

	offerParser := func(el *colly.HTMLElement) {
		el.ForEach("#search-results > li", func(_ int, e *colly.HTMLElement) {
			var offer corporation.Offer
			var err error

			// set corporation name
			offer.CorporationName = c.corporation.Name

			// get offer url
			offer.URL = e.ChildAttr("div.search-result-button a", "href")

			// skip houses under reservation
			if e.ChildText("figure span.status-sticker") == offerReserved {
				return
			}

			// get image url
			offer.RawPictureURL, err = c.parsePictureURL(e.ChildAttr("figure img.img-fluid", "src"))
			if err != nil {
				c.logger.Info("failed parsing picture url", zap.Error(err), logConnector)
			}

			// get location
			rawCity := e.ChildText("h3 > .postal-code.plaats")
			offer.Housing.CityName, err = c.parseCity(rawCity)
			if err != nil {
				return
			}
			offer.Housing.Address = fmt.Sprintf("%s, %s", e.ChildText("h3 > .street-name.straat"), rawCity)
			offer.Housing.CityDistrict, err = c.mapboxClient.CityDistrictFromAddress(offer.Housing.Address)
			if err != nil {
				c.logger.Info("could not get city district", zap.String("url", offer.URL), zap.Error(err), logConnector)
			}

			// parse price
			priceStr := strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(e.ChildText("div.huurprijs span.page-price"), ".", ""), "€", ""))
			if offer.Housing.Price, err = strconv.ParseFloat(priceStr, 32); err != nil {
				c.logger.Info("error while parsing price", zap.String("url", offer.URL), zap.Error(err), logConnector)
				return
			}

			// get housing type
			offer.Housing.Type = c.parseHousingType(e.ChildText("ul.search-result-specs > li.soortobject"))
			if offer.Housing.Type == corporation.HousingTypeUndefined {
				// skip all undefined "houses" with a price lower than 200€
				if offer.Housing.Price < 200 {
					return
				}

				// we are rewriting everything to appartement as ikwilhuren.nu does not always contains the correct housing type
				offer.Housing.Type = corporation.HousingTypeAppartement
			}

			// set minimum income (4x rent)
			offer.MinimumIncome = 12 * 4 * int(offer.Housing.Price)

			// create new offer
			offers[offer.URL] = &offer

			// visit offer url
			if err := detailCollector.Visit(offer.URL); err != nil {
				c.logger.Warn("error while checking offer details", zap.String("url", offer.URL), zap.Error(err), logConnector)
			}
		})
	}

	// add housing details
	detailCollector.OnHTML("html", func(e *colly.HTMLElement) {
		// find offer
		offerURL := e.Request.URL.String()
		if _, ok := offers[offerURL]; !ok {
			// no offer matching, no details
			return
		}

		if err := c.getHousingDetails(offers[offerURL], e); err != nil {
			c.logger.Info("error while getting house details", zap.String("address", offers[offerURL].Housing.Address), zap.Error(err), logConnector)
			return
		}

		ch <- *offers[offerURL]
	})

	// add offer
	c.collector.OnHTML("#main", offerParser) // parses first page
	paginationCollector.OnHTML("#main", offerParser)

	// parse offers
	if err := c.collector.Visit(fmt.Sprintf("%s/?action=epl_search&post_type=rental", c.corporation.URL)); err != nil {
		return err
	}

	return nil
}

func (c *client) getHousingDetails(offer *corporation.Offer, e *colly.HTMLElement) error {
	var err error

	// parse externalID
	offer.ExternalID, err = c.parseExternalID(e.ChildText("script.saswp-schema-markup-output"))
	if err != nil {
		return fmt.Errorf("error while parsing external ID: %w", err)
	}

	// parse housing characteristics
	offer.Housing.NumberBedroom, err = strconv.Atoi(e.ChildText("#Main_Aantal_Slaapkamers > dd.text"))
	if err != nil {
		c.logger.Info("error parsing number bedroom", zap.String("url", offer.URL), zap.Error(err), logConnector)
	}

	offer.Housing.Size, err = strconv.ParseFloat(strings.ReplaceAll(e.ChildText("#Main_Woonopp > dd.text"), " m²", ""), 32)
	if err != nil {
		c.logger.Info("error parsing house size", zap.String("url", offer.URL), zap.Error(err), logConnector)
	}

	rawText := strings.ToLower(e.Text)
	offer.Housing.Elevator = strings.Contains(rawText, "lift")
	offer.Housing.Garage = strings.Contains(rawText, "Parkeerplaats")

	return nil
}

func (c *client) parseHousingType(houseType string) corporation.HousingType {
	houseType = strings.ToLower(houseType)

	if strings.Contains(houseType, "appartement") ||
		strings.EqualFold(houseType, "studio") {
		return corporation.HousingTypeAppartement
	}

	if strings.Contains(houseType, "woning") ||
		strings.EqualFold(houseType, "woonhuis") ||
		strings.EqualFold(houseType, "portiekflat") ||
		strings.EqualFold(houseType, "hoekwoning") ||
		strings.EqualFold(houseType, "maisonnete") ||
		strings.EqualFold(houseType, "algemeen") {
		return corporation.HousingTypeHouse
	}

	return corporation.HousingTypeUndefined
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

func (c *client) parseExternalID(rawMessage string) (string, error) {
	rawID := []struct {
		Sku string `json:"sku"`
	}{}

	if err := json.Unmarshal([]byte(rawMessage), &rawID); err != nil {
		return "", err
	}

	if len(rawID) == 0 {
		return "", errors.New("no sku found")
	}

	return rawID[0].Sku, nil
}

func (c *client) parseCity(rawCity string) (string, error) {
	// clean city data
	if strings.Contains(rawCity, " - ") {
		rawCity = strings.Split(rawCity, " - ")[0]
	}

	if strings.Contains(rawCity, " (") {
		rawCity = strings.Split(rawCity, " (")[0]
	}

	// remove postcode
	for _, reg := range []*regexp.Regexp{
		regexp.MustCompile("[0-9]{4}[A-z]{2}"),
		regexp.MustCompile("[0-9]{4}"),
	} {
		rawCity = reg.ReplaceAllString(rawCity, "")
	}

	rawCity = strings.TrimSpace(rawCity)

	if len(rawCity) < 2 {
		return "", fmt.Errorf("city too short, skipping: got %s", rawCity)
	}

	return rawCity, nil
}
