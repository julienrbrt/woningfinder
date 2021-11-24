package mapbox

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"github.com/julienrbrt/woningfinder/pkg/networking"
	"github.com/julienrbrt/woningfinder/pkg/networking/query"
	"go.uber.org/zap"
)

// response from a Mapbox geocoding request
type response struct {
	Type     string `json:"type"`
	Features []struct {
		ID        string `json:"id"`
		Type      string `json:"type"`
		Text      string `json:"text"`
		PlaceName string `json:"place_name"`
		Context   []struct {
			ID   string `json:"id"`
			Text string `json:"text"`
		} `json:"context"`
	} `json:"features"`
}

// CityDistrictFromAddress obtains the city district name of a given address
func (c *client) CityDistrictFromAddress(address string) (string, error) {
	uuid := c.buildAddressUUID(address)

	// check if district is in cache
	if district, err := c.redisClient.Get(uuid); err == nil {
		return district, nil
	}

	// make request from mapbox
	district, err := c.getCityDistrict(address)
	if err != nil {
		return district, nil
	}

	// cache district
	if err := c.redisClient.Set(uuid, district); err != nil {
		c.logger.Error("failed saving address district in redis", zap.String("address", address), zap.Error(err))
	}

	return district, nil
}

func (c *client) buildAddressUUID(address string) string {
	return base64.StdEncoding.EncodeToString([]byte(address))
}

func (c *client) getCityDistrict(search string) (string, error) {
	// add authentication token
	query := query.Query{}
	query.Add("access_token", c.apiKey)
	query.Add("country", "nl")
	query.Add("limit", "1")

	request := networking.Request{
		Path:   fmt.Sprintf("mapbox.places/%s.json", search),
		Method: http.MethodGet,
		Query:  query,
	}

	resp, err := c.networkingClient.Send(&request)
	if err != nil {
		return "", err
	}

	var response response
	err = resp.ReadJSONBody(&response)
	if err != nil {
		return "", err
	}

	var district []string
	for _, feature := range response.Features {
		for _, c := range feature.Context {
			if strings.Contains(c.ID, "neighborhood") || strings.Contains(c.ID, "locality") {
				district = append(district, strings.ToLower(c.Text))
			}
		}
	}

	if len(district) == 0 {
		return "", fmt.Errorf("error finding city district for %s: no match found", search)
	}

	return strings.Join(district, " - "), nil
}
