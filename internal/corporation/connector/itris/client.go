package itris

import (
	"fmt"

	"github.com/gocolly/colly/v2"
	"github.com/woningfinder/woningfinder/internal/corporation/connector"
	"github.com/woningfinder/woningfinder/pkg/logging"
	"github.com/woningfinder/woningfinder/pkg/mapbox"
)

type client struct {
	collector    *colly.Collector
	logger       *logging.Logger
	mapboxClient mapbox.Client
	url          string
}

// NewClient allows to connect to itris ERP
func NewClient(logger *logging.Logger, mapboxClient mapbox.Client, url string) (connector.Client, error) {
	collector, err := connector.NewWebParsingConnector(logger, "itris")
	if err != nil {
		return nil, fmt.Errorf("error creating itris connector: %w", err)
	}

	return &client{
		collector:    collector,
		logger:       logger,
		mapboxClient: mapboxClient,
		url:          url,
	}, nil
}
