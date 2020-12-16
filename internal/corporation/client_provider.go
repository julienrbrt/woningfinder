package corporation

import (
	"fmt"
)

type Provider struct {
	Corporation Corporation
	Client      Client
}

// ClientProvider permits to get the corporation's client
type ClientProvider interface {
	List() *[]Corporation
	Get(corporation Corporation) (Client, error)
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
func (c *clientProvider) List() *[]Corporation {
	var corporations []Corporation
	for _, c := range c.providers {
		corporations = append(corporations, c.Corporation)
	}

	return &corporations
}

// Get gives the client used to query a corporation
func (c *clientProvider) Get(corporation Corporation) (Client, error) {
	for _, c := range c.providers {
		if c.Corporation.Name != corporation.Name || c.Corporation.URL != corporation.URL {
			continue
		}
		return c.Client, nil
	}

	return nil, fmt.Errorf("cannot find client for corporation: %s", corporation.Name)
}
