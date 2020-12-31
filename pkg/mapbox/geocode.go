package mapbox

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/woningfinder/woningfinder/pkg/networking"
	"github.com/woningfinder/woningfinder/pkg/networking/query"
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

// CityDistrictFromCoords obtains the city district name of a given coordinites
func (c *client) CityDistrictFromCoords(latitude, longitude string) (string, error) {
	return c.getCityDistrict(fmt.Sprintf("%s,%s", longitude, latitude))
}

// CityDistrictFromAddress obtains the city district name of a given address
func (c *client) CityDistrictFromAddress(address string) (string, error) {
	return c.getCityDistrict(address)
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

	if len(response.Features) == 0 {
		return "", fmt.Errorf("error finding city district for %s: no match found", search)
	}

	for _, c := range response.Features[0].Context {
		if strings.Contains(c.ID, "neighborhood") {
			return strings.ToLower(c.Text), nil
		}
	}

	return "", fmt.Errorf("error finding city district for %s: no match found", search)
}
