package itris

import (
	"net/http/cookiejar"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/extensions"
	"github.com/woningfinder/woningfinder/internal/corporation"
	"github.com/woningfinder/woningfinder/internal/corporation/connector"
	"github.com/woningfinder/woningfinder/pkg/logging"
	"github.com/woningfinder/woningfinder/pkg/mapbox"
	"go.uber.org/zap"
)

var logConnector = zap.String("connector", "itris")

type client struct {
	collector      *colly.Collector
	logger         *logging.Logger
	mapboxClient   mapbox.Client
	corporation    corporation.Corporation
	itrisCSRFToken string
}

// Note, if we start to get blocked investigate in proxy switcher
// https://github.com/gocolly/colly/blob/v2.1.0/_examples/proxy_switcher/proxy_switcher.go

// NewClient allows to connect to itris ERP
func NewClient(logger *logging.Logger, mapboxClient mapbox.Client, corporation corporation.Corporation) (connector.Client, error) {
	c := colly.NewCollector(
		// allow revisiting url between jobs and ignore robot txt
		colly.AllowURLRevisit(),
		colly.IgnoreRobotsTxt(),
	)

	// tweak default http client
	c.WithTransport(connector.DefaultCollyHTTPTransport)

	// add cookie jar
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}
	c.SetCookieJar(jar)

	// set random desktop user agent
	extensions.RandomUserAgent(c)

	// set limit rules
	c.Limit(&colly.LimitRule{
		RandomDelay: 2 * time.Second, // add a random delay of maximum two seconds between requests
	})

	// before making a request print the following
	c.OnRequest(func(r *colly.Request) {
		// set accept header
		r.Headers.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")

		logger.Info("visiting", zap.String("url", r.URL.String()), logConnector)
	})

	return &client{
		collector:    c,
		logger:       logger,
		mapboxClient: mapboxClient,
		corporation:  corporation,
	}, nil
}
