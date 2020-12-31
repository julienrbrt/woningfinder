package bootstrap

import (
	"github.com/woningfinder/woningfinder/internal/corporation"
	"github.com/woningfinder/woningfinder/internal/corporation/connector/itris"
	"github.com/woningfinder/woningfinder/internal/corporation/onshuis"
	"github.com/woningfinder/woningfinder/pkg/mapbox"
	"go.uber.org/zap"
)

// CreateOnsHuisClient creates a client for OnsHuis
func CreateOnsHuisClient(logger *zap.Logger, mapboxClient mapbox.Client) corporation.Client {
	return itris.NewClient(logger, mapboxClient, onshuis.Info.APIEndpoint.String())
}
