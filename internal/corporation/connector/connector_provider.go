package connector

import (
	"fmt"

	"github.com/julienrbrt/woningfinder/internal/corporation"
	"github.com/julienrbrt/woningfinder/internal/corporation/city"
)

// Provider provides all housing corporation data (client, cities, ...)
type Provider struct {
	Corporation corporation.Corporation
	Connector   Client
}

// ConnectorProvider permits to get the corporations data
type ConnectorProvider interface {
	GetClient(name string) (Client, error)
	GetCities() map[string]city.City
	GetCorporations() []corporation.Corporation
	GetCorporation(name string) (corporation.Corporation, error)
}

type connectorProvider struct {
	providers []Provider
}

func NewConnectorProvider(providers []Provider) ConnectorProvider {
	return &connectorProvider{
		providers: providers,
	}
}

// Get gives the client used to query a corporation
func (c *connectorProvider) GetClient(name string) (Client, error) {
	for _, c := range c.providers {
		if c.Corporation.Name != name {
			continue
		}

		return c.Connector, nil
	}

	return nil, fmt.Errorf("cannot find client for corporation: %s", name)
}

// GetCorporations all the supported corporations
func (c *connectorProvider) GetCorporations() []corporation.Corporation {
	var corporations []corporation.Corporation
	for _, c := range c.providers {
		corporations = append(corporations, c.Corporation)
	}

	return corporations
}

func (c *connectorProvider) GetCorporation(name string) (corporation.Corporation, error) {
	for _, c := range c.providers {
		if c.Corporation.Name != name {
			continue
		}

		return c.Corporation, nil
	}

	return corporation.Corporation{}, fmt.Errorf("cannot find corporation: %s", name)
}

func (c *connectorProvider) GetCities() map[string]city.City {
	cities := make(map[string]city.City)
	for _, corporation := range c.GetCorporations() {
		for _, city := range corporation.Cities {
			cities[city.Name] = city
		}
	}

	return cities
}
