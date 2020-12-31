package mapbox

import (
	"net/url"

	"github.com/woningfinder/woningfinder/pkg/networking"
)

// GeocodingAPIEndpoint is the geocoding mapbox api endpoint
var GeocodingAPIEndpoint = url.URL{Scheme: "https", Host: "api.mapbox.com", Path: "/geocoding/v5/"}

// Client for Mapbox
type Client interface {
	CityDistrictFromCoords(latitude, longitude string) (string, error)
	CityDistrictFromAddress(address string) (string, error)
}

type client struct {
	networkingClient networking.Client
	apiKey           string
}

// NewClient creates a client for Mapbox
func NewClient(c networking.Client, apiKey string) Client {
	return &client{
		networkingClient: c,
		apiKey:           apiKey,
	}
}
