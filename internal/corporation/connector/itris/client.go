package itris

import (
	"net/http/cookiejar"

	"go.uber.org/zap"

	"github.com/gocolly/colly/v2"
	"github.com/woningfinder/woningfinder/internal/corporation"
	"github.com/woningfinder/woningfinder/pkg/mapbox"
)

type client struct {
	collector    *colly.Collector
	logger       *zap.Logger
	mapboxClient mapbox.Client
	url          string
}

// NewClient allows to connect to itris ERP
func NewClient(logger *zap.Logger, mapboxClient mapbox.Client, url string) (corporation.Client, error) {
	c := colly.NewCollector()
	// allow revisiting url between jobs
	c.AllowURLRevisit = true

	// add cookie jar
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}
	c.SetCookieJar(jar)

	// before making a request print the following
	c.OnRequest(func(r *colly.Request) {
		logger.Sugar().Infof("itris client visiting %s", r.URL.String())
	})

	return &client{
		collector:    c,
		logger:       logger,
		mapboxClient: mapboxClient,
		url:          url,
	}, nil
}
