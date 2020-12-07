package osm

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// reverseGeo stores the reverse geocoding of a request to osm
type reverseGeo struct {
	Lat         string `json:"lat"`
	Lon         string `json:"lon"`
	Category    string `json:"category"`
	Type        string `json:"type"`
	Addresstype string `json:"addresstype"`
	DisplayName string `json:"display_name"`
	Address     struct {
		HouseNumber string `json:"house_number"`
		Road        string `json:"road"`
		Residential string `json:"residential"`
		City        string `json:"city"`
		Postcode    string `json:"postcode"`
	} `json:"address"`
}

// GetResidential obtains the residential name of a given latitude and longitude
func GetResidential(latitude, longitude string) (string, error) {
	rawURL := fmt.Sprintf("https://nominatim.openstreetmap.org/reverse?format=jsonv2&lat=%s&lon=%s", latitude, longitude)

	resp, err := http.Get(rawURL)
	if err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return "", err
	}

	var reverseGeo reverseGeo
	if err := json.Unmarshal(body, &reverseGeo); err != nil {
		return "", err
	}

	return strings.ToLower(reverseGeo.Address.Residential), nil
}
