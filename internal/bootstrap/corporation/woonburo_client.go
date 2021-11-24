package corporation

import (
	"net/http"
	"net/http/cookiejar"
	"time"

	"github.com/julienrbrt/woningfinder/internal/corporation"
	"github.com/julienrbrt/woningfinder/internal/corporation/connector"
	"github.com/julienrbrt/woningfinder/internal/corporation/connector/woonburo"
	"github.com/julienrbrt/woningfinder/pkg/logging"
	"github.com/julienrbrt/woningfinder/pkg/mapbox"
	"github.com/julienrbrt/woningfinder/pkg/networking"
	"github.com/julienrbrt/woningfinder/pkg/networking/middleware"
	"github.com/julienrbrt/woningfinder/pkg/networking/retry"
	"go.uber.org/zap"
)

// CreateWoonburoClient creates a client for Woonburo
func CreateWoonburoClient(logger *logging.Logger, mapboxClient mapbox.Client, corporation corporation.Corporation) connector.Client {
	// add cookie jar
	jar, err := cookiejar.New(nil)
	if err != nil {
		logger.Fatal("error when creating cookie-jar", zap.Error(err))
	}

	client := &http.Client{
		Jar: jar,
	}
	defaultMiddleWare := []networking.ClientMiddleware{
		middleware.CreateHostMiddleware(corporation.APIEndpoint),
		middleware.CreateDefaultHeadersMiddleware(map[string]string{
			// note if detected than blocked by user-agent check https://techblog.willshouse.com/2012/01/03/most-common-user-agents/
			"User-Agent":   "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.141 Safari/537.36",
			"Content-Type": "application/json",
		}),
		middleware.CreateRetryMiddleware(retry.DefaultRetryPolicy(), time.Sleep),
		middleware.CreateTimeoutMiddleware(middleware.DefaultRequestTimeout),
	}

	httpClient := networking.NewClient(client, defaultMiddleWare...)

	return woonburo.NewClient(logger, httpClient, mapboxClient, corporation)
}
