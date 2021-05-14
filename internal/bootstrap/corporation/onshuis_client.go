package corporation

import (
	"github.com/woningfinder/woningfinder/internal/corporation/connector"
	"github.com/woningfinder/woningfinder/internal/corporation/connector/itris"
	"github.com/woningfinder/woningfinder/pkg/logging"
	"github.com/woningfinder/woningfinder/pkg/mapbox"
)

// CreateOnsHuisClient creates a client for OnsHuis
func CreateOnsHuisClient(logger *logging.Logger, mapboxClient mapbox.Client) connector.Client {
	client, err := itris.NewClient(logger, mapboxClient, itris.OnsHuisInfo.APIEndpoint.String())
	if err != nil {
		logger.Sugar().Fatal(err)
	}

	return client
}
