package bootstrap

import (
	"github.com/woningfinder/woningfinder/internal/corporation"
	"github.com/woningfinder/woningfinder/internal/corporation/dewoonplaats"
	"github.com/woningfinder/woningfinder/internal/corporation/onshuis"
	"go.uber.org/zap"
)

// CreateClientProvider provides the client of a corporation
func CreateClientProvider(logger *zap.Logger) corporation.ClientProvider {
	providers := []corporation.Provider{
		{
			Corporation: dewoonplaats.Info,
			Client:      CreateDeWoonplaatsClient(logger),
		},
		{
			Corporation: onshuis.Info,
			Client:      CreateOnsHuisClient(logger),
		},
	}

	return corporation.NewClientProvider(providers)
}
