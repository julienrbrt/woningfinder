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
	Infoveldkort          string `json:"infoveldKort"`
	Huurtoeslagvoorwaarde struct {
		Icon              string      `json:"icon"`
		ID                interface{} `json:"id"`
		Localizedicontext interface{} `json:"localizedIconText"`
	} `json:"huurtoeslagVoorwaarde"`
	Huurtoeslagmogelijk     interface{} `json:"huurtoeslagMogelijk"`
	Specifiekevoorzieningen []struct {
		ID string `json:"id"`
	} `json:"specifiekeVoorzieningen"`
	Inschrijvingvereistvoorreageren bool   `json:"inschrijvingVereistVoorReageren"`
	Postalcode                      string `json:"postalcode"`
	Street                          string `json:"street"`
	Housenumber                     string `json:"houseNumber"`
	Housenumberaddition             string `json:"houseNumberAddition"`
	Regio                           struct {
		Name string `json:"name"`
		ID   string `json:"id"`
	} `json:"regio"`
	Municipality struct {
		Name string `json:"name"`
		ID   string `json:"id"`
	} `json:"municipality"`
	City struct {
		Name string `json:"name"`
		ID   string `json:"id"`
	} `json:"city"`
	Quarter struct {
		Name string `json:"name"`
		ID   string `json:"id"`
	} `json:"quarter"`
	Neighborhood struct {
		Name string      `json:"name"`
		ID   interface{} `json:"id"`
	} `json:"neighborhood"`
	Dwellingtype struct {
		Categorie           string `json:"categorie"`
		Huurprijsduuractief bool   `json:"huurprijsDuurActief"`
		ID                  string `json:"id"`
		Localizedname       string `json:"localizedName"`
	} `json:"dwellingType"`
	Availablefrom       string  `json:"availableFrom"`
	NetRent             float64 `json:"netRent"`
	TotalRent           float64 `json:"totalRent"`
	Flexibelhurenactief bool    `json:"flexibelHurenActief"`
	Sellingprice        int     `json:"sellingPrice"`
	Sleepingroom        struct {
		Amountofrooms string `json:"amountOfRooms"`
		ID            string `json:"id"`
	} `json:"sleepingRoom"`
	Floor struct {
		ID            interface{} `json:"id"`
		Localizedname string      `json:"localizedName"`
	} `json:"floor"`
	Balcony     bool `json:"balcony"`
	Balconysite struct {
		Localizedname string `json:"localizedName"`
	} `json:"balconySite"`
	Constructionyear int `json:"constructionYear"`
	Model            struct {
		Modelcategorie struct {
			Icon          string `json:"icon"`
			Code          string `json:"code"`
			Toonopwebsite bool   `json:"toonOpWebsite"`
			ID            string `json:"id"`
		} `json:"modelCategorie"`
		Ishospiteren                      bool `json:"isHospiteren"`
		Advertentiesluitennaeerstereactie bool `json:"advertentieSluitenNaEersteReactie"`
		Einddatumtonen                    bool `json:"einddatumTonen"`
		Aantalreactiestonen               bool `json:"aantalReactiesTonen"`
	} `json:"model"`
	Rentbuy         string        `json:"rentBuy"`
	Publicationdate string        `json:"publicationDate"`
	Closingdate     string        `json:"closingDate"`
	Latitude        string        `json:"latitude"`
	Longitude       string        `json:"longitude"`
	Floorplans      []interface{} `json:"floorplans"`
	Pictures        []struct {
		Label string `json:"label"`
		URI   string `json:"uri"`
		Type  string `json:"type"`
	} `json:"pictures"`
	Reactieurl   string  `json:"reactieUrl"`
	Newlybuild   bool    `json:"newlyBuild"`
	Areadwelling float64 `json:"areaDwelling"`
	Actionlabel  struct {
		Localizedlabel string `json:"localizedLabel"`
	} `json:"actionLabel"`
	Actionlabelfrom           string      `json:"actionLabelFrom"`
	Actionlabeluntil          string      `json:"actionLabelUntil"`
	Actionlabelifactive       bool        `json:"actionLabelIfActive"`
	Relatiehuurinkomendata    interface{} `json:"relatieHuurInkomenData"`
	Relatiehuurinkomengroepen interface{} `json:"relatieHuurInkomenGroepen"`
	Doelgroepen               []struct {
		Icon string `json:"icon"`
		Code string `json:"code"`
		ID   string `json:"id"`
	} `json:"doelgroepen"`
	Koopvoorwaarden struct {
		Localizednaam interface{} `json:"localizedNaam"`
	} `json:"koopvoorwaarden"`
	Isextraaanbod bool `json:"isExtraAanbod"`
	Vatinclusive  bool `json:"vatInclusive"`
	Woningsoort   struct {
		ID            string `json:"id"`
		Localizednaam string `json:"localizedNaam"`
	} `json:"woningsoort"`
	Aantalmedebewoners                   int    `json:"aantalMedebewoners"`
	Isexternmodeltype                    bool   `json:"isExternModelType"`
	Iszelfstandig                        bool   `json:"isZelfstandig"`
	Urlkey                               string `json:"urlKey"`
	Availablefromdate                    string `json:"availableFromDate"`
	Verzameladvertentieid                int    `json:"verzameladvertentieID"`
	ID                                   string `json:"id"`
	Isgepubliceerd                       bool   `json:"isGepubliceerd"`
	Isingepubliceerdeverzameladvertentie bool   `json:"isInGepubliceerdeVerzameladvertentie"`
}

type offerDetails struct {
	Infoveldbewoners      string `json:"infoveldBewoners"`
	Infoveldkort          string `json:"infoveldKort"`
	Hospiterenvanaf       string `json:"hospiterenVanaf"`
	Extrainformatieurl    string `json:"extraInformatieUrl"`
	Huurtoeslagvoorwaarde struct {
		Icon              string      `json:"icon"`
		ID                interface{} `json:"id"`
		Localizednaam     string      `json:"localizedNaam"`
		Localizedicontext interface{} `json:"localizedIconText"`
	} `json:"huurtoeslagVoorwaarde"`
	Sorteergroep struct {
		Code string      `json:"code"`
		ID   interface{} `json:"id"`
	} `json:"sorteergroep"`
	Huurtoeslagmogelijk        interface{} `json:"huurtoeslagMogelijk"`
	Beschikbaartot             string      `json:"beschikbaarTot"`
	Actionlabeltoelichting     string      `json:"actionLabelToelichting"`
	Huurinkomenstabelgebruiken bool        `json:"huurinkomenstabelGebruiken"`
	Voorrangurgentie           bool        `json:"voorrangUrgentie"`
	Voorrangoverigeurgenties   bool        `json:"voorrangOverigeUrgenties"`
	Voorranghuishoudgroottemin int         `json:"voorrangHuishoudgrootteMin"`
	Voorranghuishoudgroottemax int         `json:"voorrangHuishoudgrootteMax"`
	Voorrangleeftijdmin        int         `json:"voorrangLeeftijdMin"`
	Voorrangleeftijdmax        int         `json:"voorrangLeeftijdMax"`
	Voorranggezinnenkinderen   bool        `json:"voorrangGezinnenKinderen"`
	Voorrangkernbinding        bool        `json:"voorrangKernbinding"`
	Woningvoorrangvoor         struct {
		Localizedname interface{} `json:"localizedName"`
	} `json:"woningVoorrangVoor"`
	Specifiekevoorzieningen []struct {
		Description          string `json:"description"`
		Incode               string `json:"inCode"`
		Dwellingtypecategory string `json:"dwellingTypeCategory"`
		ID                   string `json:"id"`
		Localizedname        string `json:"localizedName"`
	} `json:"specifiekeVoorzieningen"`
	Reactiondata struct {
		Mogelijkepositie         interface{} `json:"mogelijkePositie"`
		Voorlopigepositie        interface{} `json:"voorlopigePositie"`
		Kanreageren              bool        `json:"kanReageren"`
		Ispassend                bool        `json:"isPassend"`
		Redenmagnietreagerencode string      `json:"redenMagNietReagerenCode"`
		Loggedin                 bool        `json:"loggedin"`
		Action                   string      `json:"action"`
		Objecttype               string      `json:"objecttype"`
		Label                    string      `json:"label"`
		Openexternelink          bool        `json:"openExterneLink"`
		Isvrijesectorwoning      bool        `json:"isVrijeSectorWoning"`
		URL                      string      `json:"url"`
	} `json:"reactionData"`
	Isvrijesectorwoning             bool `json:"isVrijeSectorWoning"`
	Inschrijvingvereistvoorreageren bool `json:"inschrijvingVereistVoorReageren"`
	Corporation                     struct {
		Name    string `json:"name"`
		Picture struct {
			Location string `json:"location"`
		} `json:"picture"`
		Website string `json:"website"`
	} `json:"corporation"`
	Postalcode          string `json:"postalcode"`
	Street              string `json:"street"`
	Housenumber         string `json:"houseNumber"`
	Housenumberaddition string `json:"houseNumberAddition"`
	Regio               struct {
		Name interface{} `json:"name"`
	} `json:"regio"`
	Municipality struct {
		Name string `json:"name"`
	} `json:"municipality"`
	City struct {
		Name string `json:"name"`
	} `json:"city"`
	Quarter struct {
		Name               string      `json:"name"`
		Extrainformatieurl interface{} `json:"extraInformatieUrl"`
		ID                 string      `json:"id"`
	} `json:"quarter"`
	Dwellingtype struct {
		Categorie           string `json:"categorie"`
		Huurprijsduuractief bool   `json:"huurprijsDuurActief"`
		Localizedname       string `json:"localizedName"`
	} `json:"dwellingType"`
	Voorrangurgentiereden struct {
		Localizedname interface{} `json:"localizedName"`
	} `json:"voorrangUrgentieReden"`
	Availablefrom       string  `json:"availableFrom"`
	Netrent             float64 `json:"netRent"`
	Calculationrent     float64 `json:"calculationRent"`
	Totalrent           float64 `json:"totalRent"`
	Flexibelhurenactief bool    `json:"flexibelHurenActief"`
	Heatingcosts        float64 `json:"heatingCosts"`
	Additionalcosts     float64 `json:"additionalCosts"`
	Servicecosts        float64 `json:"serviceCosts"`
	Sellingprice        float64 `json:"sellingPrice"`
	Description         string  `json:"description"`
	Bestemming          struct {
	} `json:"bestemming"`
	Arealivingroom   int    `json:"areaLivingRoom"`
	Areasleepingroom string `json:"areaSleepingRoom"`
	Sleepingroom     struct {
		Amountofrooms string `json:"amountOfRooms"`
		ID            string `json:"id"`
		Localizedname string `json:"localizedName"`
	} `json:"sleepingRoom"`
	Energyindex string `json:"energyIndex"`
	Floor       struct {
		Localizedname string `json:"localizedName"`
	} `json:"floor"`
	Garden     bool `json:"garden"`
	Gardensite struct {
		Localizedname string `json:"localizedName"`
	} `json:"gardenSite"`
	Oppervlaktetuin struct {
		Localizedname interface{} `json:"localizedName"`
	} `json:"oppervlakteTuin"`
	Balcony     bool `json:"balcony"`
	Balconysite struct {
		Localizedname string `json:"localizedName"`
	} `json:"balconySite"`
	Heating struct {
		Localizedname string `json:"localizedName"`
	} `json:"heating"`
	Kitchen struct {
		Localizedname string `json:"localizedName"`
	} `json:"kitchen"`
	Attic struct {
		Localizedname string `json:"localizedName"`
	} `json:"attic"`
	Constructionyear         int `json:"constructionYear"`
	Minimumincome            int `json:"minimumIncome"`
	Maximumincome            int `json:"maximumIncome"`
	Minimumhouseholdsize     int `json:"minimumHouseholdSize"`
	Maximumhouseholdsize     int `json:"maximumHouseholdSize"`
	Minimumage               int `json:"minimumAge"`
	Maximumage               int `json:"maximumAge"`
	Inwonendekinderenminimum int `json:"inwonendeKinderenMinimum"`
	Inwonendekinderenmaximum int `json:"inwonendeKinderenMaximum"`
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
	Videos                        []interface{} `json:"videos"`
	Gebruikfotoalsheader          bool          `json:"gebruikFotoAlsHeader"`
	Remainingtimeuntilclosingdate string        `json:"remainingTimeUntilClosingDate"`
	Reactieurl                    string        `json:"reactieUrl"`
	Temporaryrent                 bool          `json:"temporaryRent"`
	Showenergycosts               bool          `json:"showEnergyCosts"`
	Newlybuild                    bool          `json:"newlyBuild"`
	Storageroom                   bool          `json:"storageRoom"`
	Energycosts                   []interface{} `json:"energyCosts"`
	Areadwelling                  float64       `json:"areaDwelling"`
	Areaperceel                   string        `json:"areaPerceel"`
	Volumedwelling                string        `json:"volumeDwelling"`
	Actionlabel                   struct {
		Localizedlabel interface{} `json:"localizedLabel"`
	} `json:"actionLabel"`
	Actionlabelfrom           string      `json:"actionLabelFrom"`
	Actionlabeluntil          string      `json:"actionLabelUntil"`
	Actionlabelifactive       bool        `json:"actionLabelIfActive"`
	Relatiehuurinkomendata    interface{} `json:"relatieHuurInkomenData"`
	Relatiehuurinkomengroepen interface{} `json:"relatieHuurInkomenGroepen"`
	Doelgroepen               []struct {
		Icon string `json:"icon"`
		Code string `json:"code"`
	} `json:"doelgroepen"`
	Koopvoorwaarden struct {
		Localizednaam interface{} `json:"localizedNaam"`
	} `json:"koopvoorwaarden"`
	Koopprijstype struct {
		Localizednaam interface{} `json:"localizedNaam"`
	} `json:"koopprijsType"`
	Koopkorting struct {
		Localizednaam interface{} `json:"localizedNaam"`
	} `json:"koopkorting"`
	Koopproducten struct {
		URL           interface{} `json:"url"`
		Picture       interface{} `json:"picture"`
		Localizednaam interface{} `json:"localizedNaam"`
	} `json:"koopproducten"`
	Isextraaanbod bool          `json:"isExtraAanbod"`
	Makelaars     []interface{} `json:"makelaars"`
	Lengte        string        `json:"lengte"`
	Breedte       string        `json:"breedte"`
	Hoogte        string        `json:"hoogte"`
	Rentduration  struct {
		Incode string `json:"inCode"`
		ID     string `json:"id"`
	} `json:"rentDuration"`
	Vatinclusive                        bool `json:"vatInclusive"`
	Isgepubliceerdineenmodelmetreageren bool `json:"isGepubliceerdInEenModelMetReageren"`
	Woningsoort                         struct {
		Iszelfstandig bool   `json:"isZelfstandig"`
		ID            string `json:"id"`
		Localizednaam string `json:"localizedNaam"`
	} `json:"woningsoort"`
	Aantalmedebewoners                    int           `json:"aantalMedebewoners"`
	Isexternmodeltype                     bool          `json:"isExternModelType"`
	Iszelfstandig                         bool          `json:"isZelfstandig"`
	Urlkey                                string        `json:"urlKey"`
	Servicecomponentenbinnenservicekosten []interface{} `json:"servicecomponentenBinnenServicekosten"`
	Servicecomponentenbuitenservicekosten []interface{} `json:"servicecomponentenBuitenServicekosten"`
	Eenmaligekosten                       int           `json:"eenmaligeKosten"`
	Reactiebeleidsregels                  []interface{} `json:"reactieBeleidsregels"`
	Sorteringbeleidsregels                []interface{} `json:"sorteringBeleidsregels"`
	Complex                               struct {
		Nummer                       interface{} `json:"nummer"`
		Naam                         interface{} `json:"naam"`
		URL                          interface{} `json:"url"`
		Serviceovereenkomstverplicht interface{} `json:"serviceovereenkomstVerplicht"`
		Serviceovereenkomstkosten    interface{} `json:"serviceovereenkomstKosten"`
		ID                           interface{} `json:"id"`
	} `json:"complex"`
	Serviceovereenkomstkosten       string `json:"serviceovereenkomstKosten"`
	Extrainschrijfduuruitgeschakeld bool   `json:"extraInschrijfduurUitgeschakeld"`
	Eigenaar                        struct {
		Name    interface{} `json:"name"`
		Website interface{} `json:"website"`
		Logo    interface{} `json:"logo"`
	} `json:"eigenaar"`
	Availablefromdate                        string `json:"availableFromDate"`
	Zonnepanelen                             bool   `json:"zonnepanelen"`
	Gaslozewoning                            bool   `json:"gaslozeWoning"`
	Nulopdemeterwoning                       bool   `json:"nulOpDeMeterWoning"`
	Verzameladvertentieid                    int    `json:"verzameladvertentieID"`
	Ophaleninkomenviamijnoverheidbijreageren bool   `json:"ophalenInkomenViaMijnOverheidBijReageren"`
	ID                                       string `json:"id"`
	Isgepubliceerd                           bool   `json:"isGepubliceerd"`
	Isingepubliceerdeverzameladvertentie     bool   `json:"isInGepubliceerdeVerzameladvertentie"`
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
		if !supportedHousing(offer) && houseType == corporation.HousingTypeUndefined {
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
		numberBedroom = 1
	}

	house := corporation.Housing{
		Type:          houseType,
		Address:       fmt.Sprintf("%s %s-%s %s %s", offer.Street, offer.Housenumber, offer.Housenumberaddition, offer.Postalcode, offer.City.Name),
		CityName:      offer.City.Name,
		NumberBedroom: numberBedroom,
		Size:          offer.Areadwelling,
		Price:         offer.Totalrent,
		BuildingYear:  offer.Constructionyear,
		Garden:        offer.Garden,
		Garage:        false,
		Elevator:      true,
		Balcony:       offer.Balcony,
		Attic:         len(offer.Attic.Localizedname) > 0,
		Accessible:    false,
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
	}
}

// supportedHousing filters the offers to only houses supported by WoningFinder (no shared room for instance)
func supportedHousing(offer offerList) bool {
	if offer.Rentbuy != "Huur" {
		return false
	}

	if offer.TotalRent == 0 {
		return false
	}

	return true
}

func (c *client) parseHousingType(offer offerList) corporation.HousingType {
	if offer.Dwellingtype.Categorie != "woning" || !offer.Iszelfstandig {
		return corporation.HousingTypeUndefined
	}

	if strings.EqualFold(offer.Dwellingtype.Localizedname, "appartement") ||
		strings.Contains(strings.ToLower(offer.Dwellingtype.Localizedname), "studio") {
		return corporation.HousingTypeAppartement
	}

	return corporation.HousingTypeUndefined
}

func (c *client) parseSelectionMethod(offer *offerDetails) corporation.SelectionMethod {
	if offer.Model.Modelcategorie.Code == "inschrijfduur" {
		return corporation.SelectionRegistrationDate
	}

	if offer.Model.Modelcategorie.Code == "reactiedatum" {
		return corporation.SelectionRandom
	}

	return corporation.SelectionRandom
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
