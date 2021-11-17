package connector

import (
	"fmt"

	"github.com/woningfinder/woningfinder/internal/corporation"
)

// Provider provides a client to an housing corporation
type Provider struct {
	Corporation corporation.Corporation
	Connector   Client
}

// ClientProvider permits to get the corporation's client
type ClientProvider interface {
	Get(name string) (Client, error)
	GetAllCorporation() []corporation.Corporation
	GetCorporation(name string) (corporation.Corporation, error)
}

type clientProvider struct {
	providers []Provider
}

// NewClientProvider permits to create the mapping of a corporation to its client
func NewClientProvider(providers []Provider) ClientProvider {
	return &clientProvider{
		providers: providers,
	}
}

// Get gives the client used to query a corporation
func (c *clientProvider) Get(name string) (Client, error) {
	for _, c := range c.providers {
		if c.Corporation.Name != name {
			continue
		}

		return c.Connector, nil
	}

	return nil, fmt.Errorf("cannot find client for corporation: %s", name)
}

// GetAllCorporation all the supported corporations
func (c *clientProvider) GetAllCorporation() []corporation.Corporation {
	var corporations []corporation.Corporation
	for _, c := range c.providers {
		corporations = append(corporations, c.Corporation)
	}

	return corporations
}

func (c *clientProvider) GetCorporation(name string) (corporation.Corporation, error) {
	for _, c := range c.providers {
		if c.Corporation.Name != name {
			continue
		}

		return c.Corporation, nil
	}

	return corporation.Corporation{}, fmt.Errorf("cannot find corporation: %s", name)
}
