package bootstrap

import (
	"net/http"
	"time"

	"github.com/woningfinder/woningfinder/pkg/config"
	"github.com/woningfinder/woningfinder/pkg/mapbox"
	"github.com/woningfinder/woningfinder/pkg/networking"
	"github.com/woningfinder/woningfinder/pkg/networking/middleware"
)

// CreateMapboxClient creates a Mapbox client
func CreateMapboxClient() mapbox.Client {
	client := &http.Client{Timeout: 5 * time.Second}
	defaultMiddleWare := []networking.ClientMiddleware{
		middleware.CreateHostMiddleware(&mapbox.APIEndpoint),
	}

	httpClient := networking.NewClient(client, defaultMiddleWare...)

	return mapbox.NewClient(httpClient, config.MustGetString("MAPBOX_API_KEY"))
}
