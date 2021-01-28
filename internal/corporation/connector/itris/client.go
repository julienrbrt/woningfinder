package itris

import (
	"fmt"

	"github.com/gocolly/colly/v2"
	"github.com/woningfinder/woningfinder/internal/corporation"
	"github.com/woningfinder/woningfinder/internal/corporation/connector"
	"github.com/woningfinder/woningfinder/pkg/logging"
	"github.com/woningfinder/woningfinder/pkg/mapbox"
)

const name = "itris"

type client struct {
	collector    *colly.Collector
	logger       *logging.Logger
	mapboxClient mapbox.Client
	url          string
}

// NewClient allows to connect to itris ERP
func NewClient(logger *logging.Logger, mapboxClient mapbox.Client, url string) (corporation.Client, error) {
	collector, err := connector.NewCollyConnector(logger, name)
	if err != nil {
		return nil, fmt.Errorf("error creating %s connector: %w", name, err)
	}

	return &client{
		collector:    collector,
		logger:       logger,
		mapboxClient: mapboxClient,
		url:          url,
	}, nil
}
