package bootstrap

import (
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"

	"github.com/woningfinder/woningfinder/internal/corporation"
	"github.com/woningfinder/woningfinder/internal/corporation/city"
	"github.com/woningfinder/woningfinder/internal/corporation/connector"
	"github.com/woningfinder/woningfinder/internal/corporation/connector/dewoonplaats"
	"github.com/woningfinder/woningfinder/internal/corporation/scheduler"
	"github.com/woningfinder/woningfinder/pkg/logging"
	"github.com/woningfinder/woningfinder/pkg/mapbox"
	"github.com/woningfinder/woningfinder/pkg/networking"
	"github.com/woningfinder/woningfinder/pkg/networking/middleware"
	"github.com/woningfinder/woningfinder/pkg/networking/retry"
)

var dewoonplaatsInfo = corporation.Corporation{
	APIEndpoint: &url.URL{Scheme: "https", Host: "www.dewoonplaats.nl", Path: "/wh_services"},
	Name:        "De Woonplaats",
	URL:         "https://dewoonplaats.nl",
	Cities: []corporation.City{
		city.Enschede,
		city.Zwolle,
		city.Dinxperlo,
		city.Winterswijk,
		city.Neede,
		city.Wehl,
		city.Aalten,
		city.Groenlo,
		city.Bussum,
		city.Bredevoort,
		city.Ulft,
	},
	SelectionMethod: []corporation.SelectionMethod{
		corporation.SelectionRandom,
		corporation.SelectionFirstComeFirstServed,
	},
	SelectionTime: scheduler.CreateSelectionTime(18, 0, 0),
}

// CreateDeWoonplaatsClient creates a client for De Woonplaats
func CreateDeWoonplaatsClient(logger *logging.Logger, mapboxClient mapbox.Client) connector.Client {
	// add cookie jar
	jar, err := cookiejar.New(nil)
	if err != nil {
		logger.Sugar().Fatal(err)
	}

	client := &http.Client{
		Timeout: retry.DefaultTimeout,
		Jar:     jar,
	}
	defaultMiddleWare := []networking.ClientMiddleware{
		middleware.CreateHostMiddleware(dewoonplaatsInfo.APIEndpoint),
		middleware.CreateDefaultHeadersMiddleware(map[string]string{
			// note if detected than blocked by user-agent check https://techblog.willshouse.com/2012/01/03/most-common-user-agents/
			"User-Agent":   "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.141 Safari/537.36",
			"Content-Type": "application/json",
		}),
		middleware.CreateRetryMiddleware(retry.DefaultRetryPolicy(), time.Sleep),
		middleware.CreateTimeoutMiddleware(retry.DefaultTimeout),
	}

	httpClient := networking.NewClient(client, defaultMiddleWare...)

	return dewoonplaats.NewClient(logger, httpClient, mapboxClient)
}
