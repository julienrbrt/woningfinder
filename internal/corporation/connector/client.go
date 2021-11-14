package connector

import (
	"net/http"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/woningfinder/woningfinder/internal/corporation"
)

// Client defines a housing corporation client
type Client interface {
	Login(username, password string) error
	FetchOffers(chan<- corporation.Offer) error
	React(offer corporation.Offer) error
}

var DefaultCollyHTTPTransport = &http.Transport{
	Proxy:           http.ProxyFromEnvironment,
	MaxConnsPerHost: 50, // https://stackoverflow.com/questions/37774624/go-http-get-concurrency-and-connection-reset-by-peer
}

var DefaultCollyLimitRules = &colly.LimitRule{
	RandomDelay: 5 * time.Second, // add a random delay of maximum 5 seconds between requests
	Parallelism: 50,              // maximum parallel request
}
