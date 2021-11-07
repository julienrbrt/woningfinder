package corporation

import (
	"net/http"
	"net/http/cookiejar"
	"time"

	"github.com/woningfinder/woningfinder/internal/corporation/connector"
	"github.com/woningfinder/woningfinder/internal/corporation/connector/ikwilhuren"
	"github.com/woningfinder/woningfinder/pkg/logging"
	"github.com/woningfinder/woningfinder/pkg/mapbox"
	"github.com/woningfinder/woningfinder/pkg/networking"
	"github.com/woningfinder/woningfinder/pkg/networking/middleware"
	"github.com/woningfinder/woningfinder/pkg/networking/retry"
	"go.uber.org/zap"
)

// CreateIkWilHurenClient creates a client for ikwilhuren.nu
func CreateIkWilHurenClient(logger *logging.Logger, mapboxClient mapbox.Client) connector.Client {
	// add cookie jar
	jar, err := cookiejar.New(nil)
	if err != nil {
		logger.Fatal("error when creating cookie-jar", zap.Error(err))
	}

	client := &http.Client{
		Jar: jar,
	}
	defaultMiddleWare := []networking.ClientMiddleware{
		middleware.CreateHostMiddleware(ikwilhuren.Info.APIEndpoint),
		middleware.CreateDefaultHeadersMiddleware(map[string]string{
			// note if detected than blocked by user-agent check https://techblog.willshouse.com/2012/01/03/most-common-user-agents/
			"User-Agent":   "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.141 Safari/537.36",
			"Content-Type": "application/x-www-form-urlencoded",
		}),
		middleware.CreateRetryMiddleware(retry.DefaultRetryPolicy(), time.Sleep),
		middleware.CreateTimeoutMiddleware(middleware.DefaultTimeout),
	}

	httpClient := networking.NewClient(client, defaultMiddleWare...)

	return ikwilhuren.NewClient(logger, httpClient, mapboxClient)
}
