package connector

import (
	"net"
	"net/http"
	"time"

	"github.com/woningfinder/woningfinder/internal/corporation"
)

// Client defines a housing corporation client
type Client interface {
	Login(username, password string) error
	GetOffers() ([]corporation.Offer, error)
	React(offer corporation.Offer) error
}

var DefaultCollyHTTPTransport = &http.Transport{
	Proxy: http.ProxyFromEnvironment,
	DialContext: (&net.Dialer{
		Timeout:   10 * time.Second,
		KeepAlive: 30 * time.Second,
	}).DialContext,
	MaxIdleConns:          100,
	IdleConnTimeout:       90 * time.Second,
	TLSHandshakeTimeout:   10 * time.Second,
	ExpectContinueTimeout: 1 * time.Second,
}
