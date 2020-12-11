package bootstrap

import (
	"fmt"

	"github.com/woningfinder/woningfinder/internal/corporation"
)

// ClientProvider permits to get the corporation's client
type ClientProvider interface {
	Get(corporation *corporation.Corporation) (corporation.Client, error)
}

type clientProvider struct {
}

func NewClientProvider() ClientProvider {
	return &clientProvider{}
}

// TODO
// Get gives the client used to query a corporation
func (c *clientProvider) Get(corporation *corporation.Corporation) (corporation.Client, error) {
	// clientName := fmt.Sprintf("Create%sClient", strings.Replace(corporation.Name, " ", "", -1))
	// maybe use reflect?

	switch corporation.Name {
	case "De Woonplaats":
		return CreateDeWoonplaatsClient(), nil
	case "OnsHuis":
		return CreateOnsHuisClient(), nil
	default:
		return nil, fmt.Errorf("cannot find client for corporation: %s", corporation.Name)
	}
}
