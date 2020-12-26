package itris

import (
	"net/http/cookiejar"

	"go.uber.org/zap"

	"github.com/gocolly/colly/v2"
	"github.com/woningfinder/woningfinder/pkg/connector"
)

type itrisConnector struct {
	logger    *zap.Logger
	url       string
	collector *colly.Collector
}

func NewConnector(logger *zap.Logger, url string) connector.Connector {
	c := colly.NewCollector()
	// allow revisiting url between jobs
	c.AllowURLRevisit = true

	// add cookie jar
	jar, err := cookiejar.New(nil)
	if err != nil {
		logger.Sugar().Fatal(err)
	}
	c.SetCookieJar(jar)

	// before making a request print the following
	c.OnRequest(func(r *colly.Request) {
		logger.Sugar().Infof("itris connector visiting", r.URL.String())
	})

	return &itrisConnector{
		logger:    logger,
		url:       url,
		collector: c,
	}
}
