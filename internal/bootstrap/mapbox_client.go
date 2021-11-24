package bootstrap

import (
	"net/http"
	"time"

	"github.com/julienrbrt/woningfinder/internal/database"
	"github.com/julienrbrt/woningfinder/pkg/config"
	"github.com/julienrbrt/woningfinder/pkg/logging"
	"github.com/julienrbrt/woningfinder/pkg/mapbox"
	"github.com/julienrbrt/woningfinder/pkg/networking"
	"github.com/julienrbrt/woningfinder/pkg/networking/middleware"
	"github.com/julienrbrt/woningfinder/pkg/networking/retry"
)

// CreateMapboxClient creates a Mapbox client
func CreateMapboxClient(logger *logging.Logger, redisClient database.RedisClient) mapbox.Client {
	defaultMiddleWare := []networking.ClientMiddleware{
		middleware.CreateHostMiddleware(&mapbox.APIEndpoint),
		middleware.CreateDefaultHeadersMiddleware(map[string]string{"Content-Type": "application/json"}),
		middleware.CreateRetryMiddleware(retry.DefaultRetryPolicy(), time.Sleep),
		middleware.CreateTimeoutMiddleware(middleware.DefaultRequestTimeout),
	}

	httpClient := networking.NewClient(http.DefaultClient, defaultMiddleWare...)

	return mapbox.NewClient(logger, httpClient, redisClient, config.MustGetString("MAPBOX_API_KEY"))
}
