package corporation

import (
	"github.com/julienrbrt/woningfinder/internal/corporation/connector"
	"github.com/julienrbrt/woningfinder/internal/corporation/connector/domijn"
	"github.com/julienrbrt/woningfinder/pkg/logging"
	"github.com/julienrbrt/woningfinder/pkg/mapbox"
	"go.uber.org/zap"
)

// CreateDomijnClient creates a client for Domijn
func CreateDomijnClient(logger *logging.Logger, mapboxClient mapbox.Client) connector.Client {
	client, err := domijn.NewClient(logger, mapboxClient)
	if err != nil {
		logger.Fatal("error when creating domijn client", zap.Error(err))
	}

	return client
}
