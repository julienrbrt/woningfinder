package bootstrap

import (
	"net/http"
	"time"

	"github.com/woningfinder/woningfinder/pkg/config"
	"github.com/woningfinder/woningfinder/pkg/mapbox"
	"github.com/woningfinder/woningfinder/pkg/networking"
	"github.com/woningfinder/woningfinder/pkg/networking/middleware"
	"github.com/woningfinder/woningfinder/pkg/networking/retry"
)

// CreateMapboxClient creates a Mapbox client
func CreateMapboxClient() mapbox.Client {
	client := &http.Client{
		Timeout: retry.DefaultTimeout,
	}
	defaultMiddleWare := []networking.ClientMiddleware{
		middleware.CreateHostMiddleware(&mapbox.APIEndpoint),
		middleware.CreateDefaultHeadersMiddleware(map[string]string{"Content-Type": "application/json"}),
		middleware.CreateRetryMiddleware(retry.DefaultRetryPolicy(), time.Sleep),
	}

	httpClient := networking.NewClient(client, defaultMiddleWare...)

	return mapbox.NewClient(httpClient, config.MustGetString("MAPBOX_API_KEY"))
}
