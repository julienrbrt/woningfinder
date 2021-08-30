package bootstrap

import (
	"fmt"
	"net/http"
	"time"

	"github.com/woningfinder/woningfinder/pkg/config"
	"github.com/woningfinder/woningfinder/pkg/cryptocom"
	"github.com/woningfinder/woningfinder/pkg/networking"
	"github.com/woningfinder/woningfinder/pkg/networking/middleware"
	"github.com/woningfinder/woningfinder/pkg/networking/retry"
)

// CreateCryptoComClient creates a Crypto.com Pay client
func CreateCryptoComClient() cryptocom.Client {
	apiKey := config.MustGetString("CRYPTOCOM_API_KEY")
	webookSigningKey := config.MustGetString("CRYPTOCOM_WEBHOOK_SIGNING_KEY")

	client := &http.Client{
		Timeout: retry.DefaultTimeout,
	}
	defaultMiddleWare := []networking.ClientMiddleware{
		middleware.CreateHostMiddleware(&cryptocom.APIEndpoint),
		middleware.CreateDefaultHeadersMiddleware(map[string]string{
			"Content-Type":  "application/json",
			"Authorization": fmt.Sprintf("Bearer %s", apiKey),
		}),
		middleware.CreateRetryMiddleware(retry.DefaultRetryPolicy(), time.Sleep),
		middleware.CreateTimeoutMiddleware(retry.DefaultTimeout),
	}

	httpClient := networking.NewClient(client, defaultMiddleWare...)

	return cryptocom.NewClient(httpClient, apiKey, webookSigningKey)
}
