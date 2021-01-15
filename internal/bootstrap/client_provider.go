package bootstrap

import (
	"github.com/woningfinder/woningfinder/internal/corporation"
	"github.com/woningfinder/woningfinder/internal/corporation/dewoonplaats"
	"github.com/woningfinder/woningfinder/internal/corporation/onshuis"
	"github.com/woningfinder/woningfinder/pkg/logging"
	"github.com/woningfinder/woningfinder/pkg/mapbox"
)

// CreateClientProvider provides the client of a corporation
func CreateClientProvider(logger *logging.Logger, mapboxClient mapbox.Client) corporation.ClientProvider {
	providers := []corporation.Provider{
		{
			Corporation: dewoonplaats.Info,
			Client:      CreateDeWoonplaatsClient(logger, mapboxClient),
		},
		{
			Corporation: onshuis.Info,
			Client:      CreateOnsHuisClient(logger, mapboxClient),
		},
	}

	return corporation.NewClientProvider(providers)
}
