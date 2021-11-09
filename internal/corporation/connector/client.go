package connector

import (
	"net"
	"net/http"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/woningfinder/woningfinder/internal/corporation"
	"github.com/woningfinder/woningfinder/pkg/networking/middleware"
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
		Timeout:   middleware.DefaultTimeout,
		KeepAlive: 30 * time.Second,
	}).DialContext,
	MaxConnsPerHost:       50, // https://stackoverflow.com/questions/37774624/go-http-get-concurrency-and-connection-reset-by-peer
	IdleConnTimeout:       90 * time.Second,
	TLSHandshakeTimeout:   10 * time.Second,
	ExpectContinueTimeout: 1 * time.Second,
}

var DefaultCollyLimitRules = &colly.LimitRule{
	RandomDelay: 5 * time.Second, // add a random delay of maximum 5 seconds between requests
	Parallelism: 50,              // maximum parallel request
}
