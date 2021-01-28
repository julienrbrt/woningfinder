package connector

import (
	"net/http/cookiejar"

	"github.com/woningfinder/woningfinder/pkg/networking/retry"

	"github.com/gocolly/colly/v2"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/woningfinder/woningfinder/pkg/logging"
)

// NewCollyConnector defines a go-colly collector, the collector is used by all the connector that does web scraping
func NewCollyConnector(logger *logging.Logger, name string) (*colly.Collector, error) {
	c := colly.NewCollector()

	// set custom networking client with retries and timeout
	retryClient := retryablehttp.NewClient()
	retryClient.HTTPClient.Timeout = retry.DefaultRetryCount
	retryClient.RetryMax = 10

	c.SetClient(retryClient.StandardClient())

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
		logger.Sugar().Debugf("%s client visiting %s", name, r.URL.String())
	})

	return c, nil
}
