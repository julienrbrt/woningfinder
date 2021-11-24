package domijn

import (
	"net/http/cookiejar"

	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/extensions"
	"github.com/julienrbrt/woningfinder/internal/corporation"
	"github.com/julienrbrt/woningfinder/internal/corporation/connector"
	"github.com/julienrbrt/woningfinder/pkg/logging"
	"github.com/julienrbrt/woningfinder/pkg/mapbox"
	"github.com/julienrbrt/woningfinder/pkg/networking/middleware"
	"go.uber.org/zap"
)

var logConnector = zap.String("connector", "domijn")

type client struct {
	collector    *colly.Collector
	logger       *logging.Logger
	mapboxClient mapbox.Client
	corporation  corporation.Corporation
}

// Note, if we start to get blocked investigate in proxy switcher
// https://github.com/gocolly/colly/blob/v2.1.0/_examples/proxy_switcher/proxy_switcher.go

// NewClient allows to connect to domijn
func NewClient(logger *logging.Logger, mapboxClient mapbox.Client) (connector.Client, error) {
	c := colly.NewCollector(
		colly.AllowURLRevisit(), // allow revisiting url between jobs
		colly.IgnoreRobotsTxt(), // ignore robots.txt
	)

	// tweak default http client
	c.WithTransport(connector.DefaultCollyHTTPTransport)

	// set limit rules
	c.Limit(connector.DefaultCollyLimitRules)

	// set request timeout
	c.SetRequestTimeout(middleware.DefaultRequestTimeout)

	// add cookie jar
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}
	c.SetCookieJar(jar)

	// set random desktop user agent
	extensions.RandomUserAgent(c)

	// before making a request print the following
	c.OnRequest(func(r *colly.Request) {
		// set content type header
		r.Headers.Set("Content-Type", "application/x-www-form-urlencoded")

		logger.Info("visiting", zap.String("url", r.URL.String()), logConnector)
	})

	return &client{
		collector:    c,
		logger:       logger,
		mapboxClient: mapboxClient,
		corporation:  Info,
	}, nil
}
