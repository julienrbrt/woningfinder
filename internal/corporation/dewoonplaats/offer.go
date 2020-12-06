package dewoonplaats

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/woningfinder/woningfinder/internal/corporation"
	"github.com/woningfinder/woningfinder/pkg/networking"
)

const methodOffer = "ZoekWoningen"

// TODO to simplify
type offerResult struct {
	Aantal  int    `json:"aantal"`
	Addhtml string `json:"addhtml"`
	Latlngs []struct {
		ID   string  `json:"id"`
		Lat  float64 `json:"lat"`
		Lng  float64 `json:"lng"`
		Type string  `json:"type"`
	} `json:"latlngs"`
	Woningen []struct {
		Aanvaarding       string `json:"aanvaarding"`
		Aanvaardingsdatum string `json:"aanvaardingsdatum"`
		Adres             string `json:"adres"`
		Adrespreview      bool   `json:"adrespreview"`
		AdresUnpublished  string `json:"adres_unpublished"`
		Advid             int    `json:"advid"`
		Balkon            bool   `json:"balkon"`
		Beschikbaarheid   string `json:"beschikbaarheid"`
		BeschikbaarheidID string `json:"beschikbaarheid_id"`
		Beschrijving      string `json:"beschrijving"`
		Bijzonderheden    string `json:"bijzonderheden"`
		BogAlgemeen       string `json:"bog_algemeen"`
		BogHuuraanbieding string `json:"bog_huuraanbieding"`
		BogLigging        string `json:"bog_ligging"`
		BogMaxoppervlak   string `json:"bog_maxoppervlak"`
		BogMinoppervlak   string `json:"bog_minoppervlak"`
		BogOppervlak      string `json:"bog_oppervlak"`
		BogOverig         string `json:"bog_overig"`
		BogParkeren       string `json:"bog_parkeren"`
		BogPrijs          string `json:"bog_prijs"`
		Bouw              string `json:"bouw"`
		Bouwjaar          int    `json:"bouwjaar"`
		Brutoprijs        string `json:"brutoprijs"`
		Criteria          struct {
			CizIndicatieValid     bool   `json:"ciz_indicatie_valid"`
			FsSettingid           int    `json:"fs_settingid"`
			KinderenValid         bool   `json:"kinderen_valid"`
			MaxGezinsgrootte      int    `json:"max_gezinsgrootte"`
			MaxGezinsgrootteValid bool   `json:"max_gezinsgrootte_valid"`
			MaxInkomen            int    `json:"max_inkomen"`
			MaxInkomenValid       bool   `json:"max_inkomen_valid"`
			MaxLeeftijd           int    `json:"max_leeftijd"`
			MaxLeeftijdValid      bool   `json:"max_leeftijd_valid"`
			MinGezinsgrootte      int    `json:"min_gezinsgrootte"`
			MinGezinsgrootteValid bool   `json:"min_gezinsgrootte_valid"`
			MinInkomen            int    `json:"min_inkomen"`
			MinInkomenValid       bool   `json:"min_inkomen_valid"`
			MinLeeftijd           int    `json:"min_leeftijd"`
			MinLeeftijdValid      bool   `json:"min_leeftijd_valid"`
			Omschrijving          string `json:"omschrijving"`
			PinActief             bool   `json:"pin_actief"`
			Volgnummer            int    `json:"volgnummer"`
		} `json:"criteria"`
		Cv                 bool   `json:"cv"`
		Datum              string `json:"datum"`
		DeOgenummer        string `json:"de_ogenummer"`
		DePublicatienummer string `json:"de_publicatienummer"`
		Doelgroepen        []struct {
			Class        string `json:"class"`
			Omschrijving string `json:"omschrijving"`
			Tag          string `json:"tag"`
		} `json:"doelgroepen"`
		Energieclass            string    `json:"energieclass"`
		Energielabel            string    `json:"energielabel"`
		Energievalue            int       `json:"energievalue"`
		Etage                   string    `json:"etage"`
		Foto                    string    `json:"foto"`
		Fotobanner              string    `json:"fotobanner"`
		Garage                  bool      `json:"garage"`
		Gemeubileerd            bool      `json:"gemeubileerd"`
		GereageerdOp            string    `json:"gereageerd_op"`
		Historic                bool      `json:"historic"`
		Huisnummer              string    `json:"huisnummer"`
		Huisnummertoevoeging    string    `json:"huisnummertoevoeging"`
		Huurdersvereniging      string    `json:"huurdersvereniging"`
		Huurdersverenigingprijs float64   `json:"huurdersverenigingprijs"`
		ID                      string    `json:"id"`
		Indeling                string    `json:"indeling"`
		Isbog                   bool      `json:"isbog"`
		Ishuur                  bool      `json:"ishuur"`
		Ishuurhoog              bool      `json:"ishuurhoog"`
		Ishuurlaag              bool      `json:"ishuurlaag"`
		Iskoop                  bool      `json:"iskoop"`
		Isobject                bool      `json:"isobject"`
		Keuken                  string    `json:"keuken"`
		Koopsom                 string    `json:"koopsom"`
		Koopvoorwaarden         string    `json:"koopvoorwaarden"`
		Lat                     float64   `json:"lat"`
		Lift                    bool      `json:"lift"`
		Ligging                 string    `json:"ligging"`
		Lng                     float64   `json:"lng"`
		Loting                  bool      `json:"loting"`
		Lotingsdatum            time.Time `json:"lotingsdatum"`
		Magreageren             bool      `json:"magreageren"`
		Makelaar                int       `json:"makelaar"`
		Mapslink                string    `json:"mapslink"`
		Maxoppervlak            string    `json:"maxoppervlak"`
		Maxreacties             int       `json:"maxreacties"`
		Minoppervlak            string    `json:"minoppervlak"`
		Nettoprijs              string    `json:"nettoprijs"`
		NietReageerbaar         string    `json:"niet_reageerbaar"`
		Opmerkingen             string    `json:"opmerkingen"`
		Overview                string    `json:"overview"`
		ParkeerplaatsHuurprijs  string    `json:"parkeerplaats_huurprijs"`
		Parkeren                string    `json:"parkeren"`
		Perdirect               bool      `json:"perdirect"`
		PinActief               bool      `json:"pin_actief"`
		Plaats                  string    `json:"plaats"`
		Pmc                     int       `json:"pmc"`
		Postcode                string    `json:"postcode"`
		Preview                 bool      `json:"preview"`
		Rayoncode               string    `json:"rayoncode"`
		Reactiedatum            string    `json:"reactiedatum"`
		Reacties                int       `json:"reacties"`
		Recreatieruimte         bool      `json:"recreatieruimte"`
		Rollatortoegankelijk    bool      `json:"rollatortoegankelijk"`
		Rolstoeltoegankelijk    bool      `json:"rolstoeltoegankelijk"`
		Servicekosten           string    `json:"servicekosten"`
		SgWijk                  string    `json:"sg_wijk"`
		Showenergie             bool      `json:"showenergie"`
		Slaapkamers             int       `json:"slaapkamers"`
		SlaapBadBg              bool      `json:"slaap_bad_bg"`
		Soort                   []string  `json:"soort"`
		Stookkosten             string    `json:"stookkosten"`
		Straat                  string    `json:"straat"`
		Streetviewlink          string    `json:"streetviewlink"`
		TehuurLuxehuur          bool      `json:"tehuur_luxehuur"`
		Thumbnail               string    `json:"thumbnail"`
		Toegankelijk            bool      `json:"toegankelijk"`
		Toeslagprijs            string    `json:"toeslagprijs"`
		Tuin                    string    `json:"tuin"`
		Tweedetoilet            bool      `json:"tweedetoilet"`
		Type                    string    `json:"type"`
		Vanafprijs              bool      `json:"vanafprijs"`
		Verbruikskosten         string    `json:"verbruikskosten"`
		Verkocht                bool      `json:"verkocht"`
		Vertrekken              []struct {
			Oppervlak string `json:"oppervlak"`
			Titel     string `json:"titel"`
		} `json:"vertrekken"`
		Verwarming      string `json:"verwarming"`
		Volgnummer      int    `json:"volgnummer"`
		Wijk            string `json:"wijk"`
		Wijkid          int    `json:"wijkid"`
		Wijzigingsdatum string `json:"wijzigingsdatum"`
		Woningtype      string `json:"woningtype"`
		Woonoppervlak   string `json:"woonoppervlak"`
		WrdID           int    `json:"wrd_id"`
		Zolder          bool   `json:"zolder"`
		Zoldertrap      string `json:"zoldertrap"`
	} `json:"woningen"`
}

func (c *client) FetchOffer() ([]corporation.Housing, error) {
	req, err := c.offerRequest()
	if err != nil {
		return nil, err
	}

	resp, err := c.Send(req)
	if err != nil {
		return nil, err
	}

	//TODO
	fmt.Println(resp)

	return nil, nil
}

func (c *client) offerRequest() (networking.Request, error) {
	req := request{
		ID:     1,
		Method: methodOffer,
		Params: []interface{}{
			struct {
				Param1 bool `json:"tehuur"`
			}{Param1: true},
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
