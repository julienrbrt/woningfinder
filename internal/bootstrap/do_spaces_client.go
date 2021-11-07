package bootstrap

import (
	"net/http"
	"time"

	"github.com/woningfinder/woningfinder/pkg/config"
	"github.com/woningfinder/woningfinder/pkg/digitalocean/spaces"
	"github.com/woningfinder/woningfinder/pkg/logging"
	"github.com/woningfinder/woningfinder/pkg/networking"
	"github.com/woningfinder/woningfinder/pkg/networking/middleware"
	"github.com/woningfinder/woningfinder/pkg/networking/retry"
	"go.uber.org/zap"
)

// CreateDOSpacesClient creates the DigitalOcean spaces client
func CreateDOSpacesClient(logger *logging.Logger) spaces.Client {
	endpoint := "ams3.digitaloceanspaces.com"
	bucketName := "woningfinder"
	accessKey := config.MustGetString("DO_SPACES_KEY")
	secretKey := config.MustGetString("DO_SPACES_SECRET")

	defaultMiddleWare := []networking.ClientMiddleware{
		middleware.CreateDefaultHeadersMiddleware(map[string]string{"Content-Type": "application/json"}),
		middleware.CreateRetryMiddleware(retry.DefaultRetryPolicy(), time.Sleep),
		middleware.CreateTimeoutMiddleware(middleware.DefaultTimeout),
	}

	httpClient := networking.NewClient(&http.Client{}, defaultMiddleWare...)
	client, err := spaces.NewClient(logger, httpClient, endpoint, bucketName, accessKey, secretKey)
	if err != nil {
		logger.Fatal("error when creating digitalocean spaces client", zap.Error(err))
	}

	return client
}
