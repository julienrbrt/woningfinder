package mapbox

import (
	"net/url"

	"github.com/julienrbrt/woningfinder/internal/database"
	"github.com/julienrbrt/woningfinder/pkg/logging"
	"github.com/julienrbrt/woningfinder/pkg/networking"
)

// APIEndpoint is the geocoding mapbox api endpoint
var APIEndpoint = url.URL{Scheme: "https", Host: "api.mapbox.com", Path: "/geocoding/v5/"}

// Client for Mapbox
type Client interface {
	CityDistrictFromAddress(address string) (string, error)
}

type client struct {
	logger           *logging.Logger
	networkingClient networking.Client
	dbClient         database.DBClient
	apiKey           string
}

// NewClient creates a client for Mapbox
func NewClient(logger *logging.Logger, c networking.Client, dbClient database.DBClient, apiKey string) Client {
	return &client{
		logger:           logger,
		networkingClient: c,
		dbClient:         dbClient,
		apiKey:           apiKey,
	}
}
