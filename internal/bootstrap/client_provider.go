package bootstrap

import (
	"github.com/woningfinder/woningfinder/internal/corporation"
	"github.com/woningfinder/woningfinder/internal/corporation/dewoonplaats"
	"github.com/woningfinder/woningfinder/internal/corporation/onshuis"
)

var providers = []corporation.Provider{
	{
		Corporation: dewoonplaats.Info,
		Client:      CreateDeWoonplaatsClient(),
	},
	{
		Corporation: onshuis.Info,
		Client:      CreateOnsHuisClient(),
	},
}

// CreateClientProvider provides the client of a corporation
func CreateClientProvider() corporation.ClientProvider {
	return corporation.NewClientProvider(providers)
}
