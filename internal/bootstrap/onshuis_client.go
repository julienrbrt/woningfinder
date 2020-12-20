package bootstrap

import (
	"github.com/woningfinder/woningfinder/internal/corporation"
	"github.com/woningfinder/woningfinder/internal/corporation/onshuis"
	"github.com/woningfinder/woningfinder/pkg/connector/itris"
)

// CreateOnsHuisClient creates a client for OnsHuis
func CreateOnsHuisClient() corporation.Client {
	return itris.NewConnector(onshuis.Info.APIEndpoint.String())
}
