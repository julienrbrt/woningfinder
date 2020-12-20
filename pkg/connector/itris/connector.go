package itris

import (
	"log"
	"net/http/cookiejar"

	"github.com/gocolly/colly/v2"
	"github.com/woningfinder/woningfinder/pkg/connector"
)

type itrisConnector struct {
	url       string
	collector *colly.Collector
}

func NewConnector(url string) connector.Connector {
	c := colly.NewCollector()

	// add cookie jar
	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatal(err)
	}
	c.SetCookieJar(jar)

	// before making a request print the following
	c.OnRequest(func(r *colly.Request) {
		log.Println("itris connector visiting", r.URL.String())
	})

	return &itrisConnector{
		url:       url,
		collector: c,
	}
}
