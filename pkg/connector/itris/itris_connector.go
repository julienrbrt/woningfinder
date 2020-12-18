package itris

import (
	"log"
	"net/url"

	"github.com/gocolly/colly"
	"github.com/woningfinder/woningfinder/pkg/connector"
)

type itrisConnector struct {
	url       url.URL
	collector *colly.Collector
}

func NewConnector(url url.URL) connector.Connector {
	c := colly.NewCollector(
		// Visit only given domain
		colly.AllowedDomains(url.Host),
	)

	// before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		log.Println("itris connector visiting", r.URL.String())
	})

	return &itrisConnector{
		url:       url,
		collector: c,
	}
}
