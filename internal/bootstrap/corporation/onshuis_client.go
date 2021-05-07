package corporation

import (
	"net/url"

	"github.com/woningfinder/woningfinder/internal/corporation"
	"github.com/woningfinder/woningfinder/internal/corporation/city"
	"github.com/woningfinder/woningfinder/internal/corporation/connector"
	"github.com/woningfinder/woningfinder/internal/corporation/connector/itris"
	"github.com/woningfinder/woningfinder/pkg/logging"
	"github.com/woningfinder/woningfinder/pkg/mapbox"
)

var OnsHuisInfo = corporation.Corporation{
	APIEndpoint: &url.URL{Scheme: "https", Host: "mijn.onshuis.com"},
	Name:        "OnsHuis",
	URL:         "https://onshuis.com",
	Cities: []corporation.City{
		city.Enschede,
		city.Hengelo,
	},
	SelectionMethod: []corporation.SelectionMethod{
		corporation.SelectionRandom,
	},
}

// CreateOnsHuisClient creates a client for OnsHuis
func CreateOnsHuisClient(logger *logging.Logger, mapboxClient mapbox.Client) connector.Client {
	client, err := itris.NewClient(logger, mapboxClient, OnsHuisInfo.APIEndpoint.String())
	if err != nil {
		logger.Sugar().Fatal(err)
	}

	return client
}
