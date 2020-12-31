package bootstrap

import (
	"net/http"
	"time"

	"github.com/woningfinder/woningfinder/pkg/config"
	"github.com/woningfinder/woningfinder/pkg/mapbox"
	"github.com/woningfinder/woningfinder/pkg/networking"
	"github.com/woningfinder/woningfinder/pkg/networking/middleware"
)

// CreateMapboxGeocodingClient creates a client geocoding Mapbox
func CreateMapboxGeocodingClient() mapbox.Client {
	client := &http.Client{Timeout: 5 * time.Second}
	defaultMiddleWare := []networking.ClientMiddleware{
		middleware.CreateHostMiddleware(&mapbox.GeocodingAPIEndpoint),
	}

	httpClient := networking.NewClient(client, defaultMiddleWare...)

	return mapbox.NewClient(httpClient, config.MustGetString("MAPBOX_API_KEY"))
}
