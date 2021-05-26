package connector

import (
	"fmt"

	"github.com/woningfinder/woningfinder/internal/corporation"
)

// Provider provides a client to an housing corporation
type Provider struct {
	Corporation corporation.Corporation
	ClientFunc  func() Client
}

// ClientProvider permits to get the corporation's client
type ClientProvider interface {
	List() []corporation.Corporation
	Get(name string) (func() Client, error)
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

// List all the supported corporations
func (c *clientProvider) List() []corporation.Corporation {
	var corporations []corporation.Corporation
	for _, c := range c.providers {
		corporations = append(corporations, c.Corporation)
	}

	return corporations
}

// Get gives the client used to query a corporation
func (c *clientProvider) Get(name string) (func() Client, error) {
	for _, c := range c.providers {
		if c.Corporation.Name != name {
			continue
		}

		return c.ClientFunc, nil
	}

	return nil, fmt.Errorf("cannot find client for corporation: %s", name)
}
