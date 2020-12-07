package bootstrap

import (
	"log"
	"net/http"
	"net/http/cookiejar"
	"time"

	"github.com/woningfinder/woningfinder/internal/corporation"
	"github.com/woningfinder/woningfinder/internal/corporation/dewoonplaats"
	"github.com/woningfinder/woningfinder/pkg/networking"
	"github.com/woningfinder/woningfinder/pkg/networking/middleware"
)

// CreateDWClient creates a client for De Woonplaats
func CreateDWClient() corporation.Client {
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
		middleware.CreateHostMiddleware(dewoonplaats.Host),
		middleware.CreateDefaultHeadersMiddleware(map[string]string{"Content-Type": "application/json"}),
	}

	httpClient := networking.NewClient(client, defaultMiddleWare...)

	return dewoonplaats.NewClient(httpClient)
}
