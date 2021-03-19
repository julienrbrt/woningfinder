package bootstrap

import (
	"github.com/woningfinder/woningfinder/internal/corporation"
	"github.com/woningfinder/woningfinder/pkg/logging"
	"github.com/woningfinder/woningfinder/pkg/mapbox"
)

// CreateClientProvider provides the client of a corporation
func CreateClientProvider(logger *logging.Logger, mapboxClient mapbox.Client) corporation.ClientProvider {
	providers := []corporation.Provider{
		{
			Corporation: dewoonplaatsInfo,
			Client:      CreateDeWoonplaatsClient(logger, mapboxClient),
		},
		{
			Corporation: onshuisInfo,
			Client:      CreateOnsHuisClient(logger, mapboxClient),
		},
	}

	return corporation.NewClientProvider(providers)
}
