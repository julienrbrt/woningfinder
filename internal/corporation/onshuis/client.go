package onshuis

import (
	"net/url"

	"github.com/woningfinder/woningfinder/internal/corporation"
	"github.com/woningfinder/woningfinder/pkg/networking"
)

// Host defines the OnsHuis offers domain
var Host = &url.URL{Scheme: "https", Host: "mijn.onshuis.com", Path: "/apps/com.itris.klantportaal/"}

type client struct {
	corporation      corporation.Corporation
	networkingClient networking.Client
}

// NewClient creates a client for a housing coporation
func NewClient(corporation corporation.Corporation, c networking.Client) corporation.Client {
	return &client{
		corporation:      corporation,
		networkingClient: c,
	}
}
