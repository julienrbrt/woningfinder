package bootstrap

import (
	"net/http"
	"time"

	"github.com/woningfinder/woningfinder/internal/database"
	"github.com/woningfinder/woningfinder/pkg/config"
	"github.com/woningfinder/woningfinder/pkg/logging"
	"github.com/woningfinder/woningfinder/pkg/mapbox"
	"github.com/woningfinder/woningfinder/pkg/networking"
	"github.com/woningfinder/woningfinder/pkg/networking/middleware"
	"github.com/woningfinder/woningfinder/pkg/networking/retry"
)

// CreateMapboxClient creates a Mapbox client
func CreateMapboxClient(logger *logging.Logger, redisClient database.RedisClient) mapbox.Client {
	defaultMiddleWare := []networking.ClientMiddleware{
		middleware.CreateHostMiddleware(&mapbox.APIEndpoint),
		middleware.CreateDefaultHeadersMiddleware(map[string]string{"Content-Type": "application/json"}),
		middleware.CreateRetryMiddleware(retry.DefaultRetryPolicy(), time.Sleep),
		middleware.CreateTimeoutMiddleware(middleware.DefaultTimeout),
	}

	httpClient := networking.NewClient(&http.Client{}, defaultMiddleWare...)

	return mapbox.NewClient(logger, httpClient, redisClient, config.MustGetString("MAPBOX_API_KEY"))
}
