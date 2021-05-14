package corporation

import (
	"github.com/woningfinder/woningfinder/internal/corporation/connector"
	"github.com/woningfinder/woningfinder/internal/corporation/connector/domijn"
	"github.com/woningfinder/woningfinder/pkg/logging"
	"github.com/woningfinder/woningfinder/pkg/mapbox"
)

// CreateDomijnClient creates a client for Domijn
func CreateDomijnClient(logger *logging.Logger, mapboxClient mapbox.Client) connector.Client {
	client, err := domijn.NewClient(logger, mapboxClient)
	if err != nil {
		logger.Sugar().Fatal(err)
	}

	return client
}
