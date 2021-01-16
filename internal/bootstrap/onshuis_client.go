package bootstrap

import (
	"github.com/woningfinder/woningfinder/internal/corporation"
	"github.com/woningfinder/woningfinder/internal/corporation/connector/itris"
	"github.com/woningfinder/woningfinder/internal/corporation/onshuis"
	"github.com/woningfinder/woningfinder/pkg/logging"
	"github.com/woningfinder/woningfinder/pkg/mapbox"
)

// CreateOnsHuisClient creates a client for OnsHuis
func CreateOnsHuisClient(logger *logging.Logger, mapboxClient mapbox.Client) corporation.Client {
	client, err := itris.NewClient(logger, mapboxClient, onshuis.Info.APIEndpoint.String())
	if err != nil {
		logger.Sugar().Fatal(err)
	}

	return client
}
