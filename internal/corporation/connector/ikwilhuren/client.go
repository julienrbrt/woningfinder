package ikwilhuren

import (
	"encoding/json"
	"net"
	"net/http"
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

var logConnector = zap.String("connector", "ikwilhuren")

type client struct {
	logger           *logging.Logger
	networkingClient networking.Client
	collector        *colly.Collector
	mapboxClient     mapbox.Client
	corporation      corporation.Corporation
}

func NewClient(logger *logging.Logger, c networking.Client, mapboxClient mapbox.Client) connector.Client {
	collector, err := getCollector(logger)
	if err != nil {
		panic(err)
	}

	return &client{
		logger:           logger,
		networkingClient: c,
		collector:        collector,
		mapboxClient:     mapboxClient,
		corporation:      Info,
	}
}

// Note, if we start to get blocked investigate in proxy switcher
// https://github.com/gocolly/colly/blob/v2.1.0/_examples/proxy_switcher/proxy_switcher.go
func getCollector(logger *logging.Logger) (*colly.Collector, error) {
	c := colly.NewCollector(
		// allow revisiting url between jobs and ignore robot txt
		colly.AllowURLRevisit(),
		colly.IgnoreRobotsTxt(),
	)

	// tweak default http client
	c.WithTransport(&http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   5 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	})

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
		// set content type header
		r.Headers.Set("Content-Type", "application/x-www-form-urlencoded")

		logger.Info("visiting", zap.String("url", r.URL.String()), logConnector)
	})

	return c, nil
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
