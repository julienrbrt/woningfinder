package connector

import (
	"net/http/cookiejar"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/extensions"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/woningfinder/woningfinder/pkg/logging"
	"github.com/woningfinder/woningfinder/pkg/networking/retry"
)

// Note, if we start to get blocked investigate in proxy switcher
// https://github.com/gocolly/colly/blob/v2.1.0/_examples/proxy_switcher/proxy_switcher.go

// NewCollyConnector defines a go-colly collector, the collector is used by all the connector that does web scraping
func NewCollyConnector(logger *logging.Logger, name string) (*colly.Collector, error) {
	c := colly.NewCollector()

	// set custom networking client with retries and timeout
	retryClient := retryablehttp.NewClient()
	retryClient.HTTPClient.Timeout = retry.DefaultTimeout
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

	// set random desktop user agent
	extensions.RandomUserAgent(c)

	// set limit rules
	c.Limit(&colly.LimitRule{
		RandomDelay: 2 * time.Second, // add a random delay of maximum two seconds between requests
	})

	// before making a request print the following
	c.OnRequest(func(r *colly.Request) {
		logger.Sugar().Debugf("%s client visiting %s", name, r.URL.String())
	})

	return c, nil
}
