package corporation

import (
	"github.com/julienrbrt/woningfinder/internal/corporation"
	"github.com/julienrbrt/woningfinder/internal/corporation/connector"
	"github.com/julienrbrt/woningfinder/internal/corporation/connector/itris"
	"github.com/julienrbrt/woningfinder/pkg/logging"
	"github.com/julienrbrt/woningfinder/pkg/mapbox"
	"go.uber.org/zap"
)

// CreateItrisClient creates a client for OnsHuis
func CreateItrisClient(logger *logging.Logger, mapboxClient mapbox.Client, corporation corporation.Corporation) connector.Client {
	client, err := itris.NewClient(logger, mapboxClient, corporation)
	if err != nil {
		logger.Fatal("error when creating itris client", zap.Error(err))
	}

	return client
}
