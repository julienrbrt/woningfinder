package bootstrap

import (
	"github.com/woningfinder/woningfinder/internal/corporation"
	"github.com/woningfinder/woningfinder/internal/corporation/onshuis"
	"github.com/woningfinder/woningfinder/pkg/connector/itris"
	"go.uber.org/zap"
)

// CreateOnsHuisClient creates a client for OnsHuis
func CreateOnsHuisClient(logger *zap.Logger) corporation.Client {
	return itris.NewConnector(logger, onshuis.Info.APIEndpoint.String())
}
