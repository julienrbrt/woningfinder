package bootstrap

import (
	"net/http"
	"time"

	"github.com/woningfinder/woningfinder/internal/config"
	"github.com/woningfinder/woningfinder/internal/networking"
	"github.com/woningfinder/woningfinder/internal/networking/middleware"
	"github.com/woningfinder/woningfinder/internal/services/mapbox"
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
