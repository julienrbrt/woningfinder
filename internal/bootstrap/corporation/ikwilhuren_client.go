package corporation

import (
	"github.com/woningfinder/woningfinder/internal/corporation/connector"
	"github.com/woningfinder/woningfinder/internal/corporation/connector/ikwilhuren"
	"github.com/woningfinder/woningfinder/pkg/logging"
	"github.com/woningfinder/woningfinder/pkg/mapbox"
)

// CreateIkWilHurenClient creates a client for ikwilhuren.nu
func CreateIkWilHurenClient(logger *logging.Logger, mapboxClient mapbox.Client) connector.Client {
	client, err := ikwilhuren.NewClient(logger, mapboxClient)
	if err != nil {
		logger.Sugar().Fatal(err)
	}

	return client
}
