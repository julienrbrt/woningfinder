package woonburo

import (
	"github.com/woningfinder/woningfinder/internal/corporation"
	"github.com/woningfinder/woningfinder/internal/corporation/connector"
	"github.com/woningfinder/woningfinder/pkg/logging"
	"github.com/woningfinder/woningfinder/pkg/mapbox"
	"github.com/woningfinder/woningfinder/pkg/networking"
)

type client struct {
	logger           *logging.Logger
	networkingClient networking.Client
	mapboxClient     mapbox.Client
	corporation      corporation.Corporation
}

// NewClient creates a client for Woonburo
func NewClient(logger *logging.Logger, c networking.Client, mapboxClient mapbox.Client, corporation corporation.Corporation) connector.Client {
	return &client{
		logger:           logger,
		networkingClient: c,
		mapboxClient:     mapboxClient,
		corporation:      corporation,
	}
}
