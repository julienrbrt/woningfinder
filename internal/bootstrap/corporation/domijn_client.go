package corporation

import (
	"net/url"

	"github.com/woningfinder/woningfinder/internal/corporation"
	"github.com/woningfinder/woningfinder/internal/corporation/city"
	"github.com/woningfinder/woningfinder/internal/corporation/connector"
	"github.com/woningfinder/woningfinder/internal/corporation/connector/domijn"
	"github.com/woningfinder/woningfinder/pkg/logging"
	"github.com/woningfinder/woningfinder/pkg/mapbox"
)

var DomijnInfo = corporation.Corporation{
	APIEndpoint: &url.URL{Scheme: "https", Host: "www.domijn.nl"},
	Name:        "Domijn",
	URL:         "https://www.domijn.nl",
	Cities: []corporation.City{
		city.Enschede,
		city.Haaksbergen,
		city.Losser,
		city.Overdinkel,
	},
	SelectionMethod: []corporation.SelectionMethod{
		corporation.SelectionRandom,
	},
}

// CreateDomijnClient creates a client for Domijn
func CreateDomijnClient(logger *logging.Logger, mapboxClient mapbox.Client) connector.Client {
	client, err := domijn.NewClient(logger, mapboxClient, DomijnInfo.APIEndpoint.String())
	if err != nil {
		logger.Sugar().Fatal(err)
	}

	return client
}
