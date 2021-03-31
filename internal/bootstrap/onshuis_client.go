package bootstrap

import (
	"net/url"

	"github.com/woningfinder/woningfinder/internal/corporation"
	"github.com/woningfinder/woningfinder/internal/corporation/connector/itris"
	"github.com/woningfinder/woningfinder/internal/domain/entity"
	"github.com/woningfinder/woningfinder/pkg/logging"
	"github.com/woningfinder/woningfinder/pkg/mapbox"
)

var onshuisInfo = entity.Corporation{
	APIEndpoint: &url.URL{Scheme: "https", Host: "mijn.onshuis.com", Path: "/apps/com.itris.klantportaal"},
	Name:        "OnsHuis",
	URL:         "https://mijn.onshuis.com",
	Cities: []entity.City{
		{Name: "Enschede"},
		{Name: "Hengelo"},
	},
	SelectionMethod: []entity.SelectionMethod{
		entity.SelectionRandom,
	},
}

// CreateOnsHuisClient creates a client for OnsHuis
func CreateOnsHuisClient(logger *logging.Logger, mapboxClient mapbox.Client) corporation.Client {
	client, err := itris.NewClient(logger, mapboxClient, onshuisInfo.APIEndpoint.String())
	if err != nil {
		logger.Sugar().Fatal(err)
	}

	return client
}
