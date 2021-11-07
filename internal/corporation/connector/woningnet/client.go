package woningnet

import (
	"encoding/json"
	"net/http/cookiejar"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/extensions"
	"github.com/woningfinder/woningfinder/internal/corporation"
	"github.com/woningfinder/woningfinder/internal/corporation/connector"
	"github.com/woningfinder/woningfinder/pkg/logging"
	"github.com/woningfinder/woningfinder/pkg/mapbox"
	"github.com/woningfinder/woningfinder/pkg/networking"
	"go.uber.org/zap"
)

var logConnector = zap.String("connector", "woningnet")

type client struct {
	logger           *logging.Logger
	networkingClient networking.Client
	collector        *colly.Collector
	mapboxClient     mapbox.Client
	corporation      corporation.Corporation
}

// NewClient creates a client for WoningNet
func NewClient(logger *logging.Logger, c networking.Client, mapboxClient mapbox.Client, corporation corporation.Corporation) connector.Client {
	collector, err := getCollector(logger)
	if err != nil {
		panic(err)
	}

	return &client{
		logger:           logger,
		networkingClient: c,
		collector:        collector,
		mapboxClient:     mapboxClient,
		corporation:      corporation,
	}
}

// Note, if we start to get blocked investigate in proxy switcher
// https://github.com/gocolly/colly/blob/v2.1.0/_examples/proxy_switcher/proxy_switcher.go
func getCollector(logger *logging.Logger) (*colly.Collector, error) {
	collector := colly.NewCollector(
		// allow revisiting url between jobs and ignore robot txt
		colly.AllowURLRevisit(),
		colly.IgnoreRobotsTxt(),
	)

	// tweak default http client
	collector.WithTransport(connector.DefaultCollyHTTPTransport)

	// add cookie jar
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}
	collector.SetCookieJar(jar)

	// set random desktop user agent
	extensions.RandomUserAgent(collector)

	// set limit rules
	collector.Limit(&colly.LimitRule{
		RandomDelay: 2 * time.Second, // add a random delay of maximum two seconds between requests
	})

	// before making a request print the following
	collector.OnRequest(func(r *colly.Request) {
		// set accept header
		r.Headers.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")

		logger.Info("visiting", zap.String("url", r.URL.String()), logConnector)
	})

	return collector, nil
}

func (c *client) Send(req networking.Request) (json.RawMessage, error) {
	// send request to networking client
	resp, err := c.networkingClient.Send(&req)
	if err != nil {
		return nil, err
	}

	var response json.RawMessage
	err = resp.ReadJSONBody(&response)
	if err != nil {
		return nil, err
	}

	return response, nil
}
