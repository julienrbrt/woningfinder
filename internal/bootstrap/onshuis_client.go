package bootstrap

import (
	"log"
	"net/http"
	"net/http/cookiejar"
	"time"

	"github.com/woningfinder/woningfinder/internal/corporation"
	"github.com/woningfinder/woningfinder/internal/corporation/onshuis"
	"github.com/woningfinder/woningfinder/pkg/networking"
	"github.com/woningfinder/woningfinder/pkg/networking/middleware"
)

var onshuisInfo = corporation.Corporation{
	Name: "OnsHuis",
	URL:  "https://mijn.onshuis.com",
	Cities: []corporation.City{
		{
			Name:   "Enschede",
			Region: "Overijssel",
		},
		{
			Name:   "Hengelo",
			Region: "Overijssel",
		},
	},
	SelectionMethod: []corporation.SelectionMethod{
		{
			Method: corporation.SelectionRandom,
		},
		{
			Method: corporation.SelectionRandom,
		},
	},
}

// CreateOnsHuisClient creates a client for OnsHuis
func CreateOnsHuisClient() corporation.Client {
	// add cookie jar
	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatal(err)
	}

	client := &http.Client{
		Timeout: 5 * time.Second,
		Jar:     jar,
	}
	defaultMiddleWare := []networking.ClientMiddleware{
		middleware.CreateHostMiddleware(onshuis.Host),
		middleware.CreateDefaultHeadersMiddleware(map[string]string{"Content-Type": "application/json"}),
	}

	httpClient := networking.NewClient(client, defaultMiddleWare...)

	return onshuis.NewClient(onshuisInfo, httpClient)
}
