package corporation

import (
	"net/http"
	"net/http/cookiejar"
	"time"

	"github.com/woningfinder/woningfinder/internal/corporation/connector"
	"github.com/woningfinder/woningfinder/internal/corporation/connector/zig"
	"github.com/woningfinder/woningfinder/pkg/logging"
	"github.com/woningfinder/woningfinder/pkg/mapbox"
	"github.com/woningfinder/woningfinder/pkg/networking"
	"github.com/woningfinder/woningfinder/pkg/networking/middleware"
	"github.com/woningfinder/woningfinder/pkg/networking/retry"
)

// CreateRoomspotClient creates a client for Roomspot
func CreateRoomspotClient(logger *logging.Logger, mapboxClient mapbox.Client) connector.Client {
	// add cookie jar
	jar, err := cookiejar.New(nil)
	if err != nil {
		logger.Sugar().Fatal(err)
	}

	client := &http.Client{
		Timeout: retry.DefaultTimeout,
		Jar:     jar,
	}
	defaultMiddleWare := []networking.ClientMiddleware{
		middleware.CreateHostMiddleware(zig.RoomspotInfo.APIEndpoint),
		middleware.CreateDefaultHeadersMiddleware(map[string]string{
			// note if detected than blocked by user-agent check https://techblog.willshouse.com/2012/01/03/most-common-user-agents/
			"User-Agent":   "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.141 Safari/537.36",
			"Content-Type": "application/json",
		}),
		middleware.CreateRetryMiddleware(retry.DefaultRetryPolicy(), time.Sleep),
		middleware.CreateTimeoutMiddleware(retry.DefaultTimeout),
	}

	httpClient := networking.NewClient(client, defaultMiddleWare...)

	return zig.NewClient(logger, httpClient, mapboxClient)
}
