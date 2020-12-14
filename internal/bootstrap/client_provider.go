package bootstrap

import (
	"github.com/woningfinder/woningfinder/internal/corporation"
)

var providers = []corporation.Provider{
	{
		Corporation: dewoonplaatsInfo,
		Client:      CreateDeWoonplaatsClient(),
	},
	{
		Corporation: onshuisInfo,
		Client:      CreateOnsHuisClient(),
	},
}

// CreateClientProvider provides the client of a corporation
func CreateClientProvider() corporation.ClientProvider {
	return corporation.NewClientProvider(providers)
}
