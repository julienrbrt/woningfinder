package bootstrap

import (
	"net/http"
	"net/http/cookiejar"
	"time"

	"github.com/woningfinder/woningfinder/internal/corporation"
	"github.com/woningfinder/woningfinder/internal/corporation/dewoonplaats"
	"github.com/woningfinder/woningfinder/pkg/networking"
	"github.com/woningfinder/woningfinder/pkg/networking/middleware"
	"go.uber.org/zap"
)

// CreateDeWoonplaatsClient creates a client for De Woonplaats
func CreateDeWoonplaatsClient(logger *zap.Logger) corporation.Client {
	// add cookie jar
	jar, err := cookiejar.New(nil)
	if err != nil {
		logger.Sugar().Fatal(err)
	}

	client := &http.Client{
		Timeout: 5 * time.Second,
		Jar:     jar,
	}
	defaultMiddleWare := []networking.ClientMiddleware{
		middleware.CreateHostMiddleware(dewoonplaats.Info.APIEndpoint),
		middleware.CreateDefaultHeadersMiddleware(map[string]string{"Content-Type": "application/json"}),
	}

	httpClient := networking.NewClient(client, defaultMiddleWare...)

	return dewoonplaats.NewClient(httpClient)
}
