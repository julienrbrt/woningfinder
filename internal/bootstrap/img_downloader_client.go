package bootstrap

import (
	"net/http"
	"time"

	"github.com/julienrbrt/woningfinder/pkg/config"
	"github.com/julienrbrt/woningfinder/pkg/downloader"
	"github.com/julienrbrt/woningfinder/pkg/logging"
	"github.com/julienrbrt/woningfinder/pkg/networking"
	"github.com/julienrbrt/woningfinder/pkg/networking/middleware"
	"github.com/julienrbrt/woningfinder/pkg/networking/retry"
	"go.uber.org/zap"
)

// CreateImgDownloader creates an image downloader
func CreateImgDownloader(logger *logging.Logger) downloader.Client {
	defaultMiddleWare := []networking.ClientMiddleware{
		middleware.CreateDefaultHeadersMiddleware(map[string]string{"Content-Type": "application/json"}),
		middleware.CreateRetryMiddleware(retry.DefaultRetryPolicy(), time.Sleep),
		middleware.CreateTimeoutMiddleware(middleware.DefaultRequestTimeout),
	}

	httpClient := networking.NewClient(http.DefaultClient, defaultMiddleWare...)
	client, err := downloader.NewClient(logger, httpClient, config.MustGetString("APP_IMG_DOWNLOADER"))
	if err != nil {
		logger.Fatal("error when creating image downloader client", zap.Error(err))
	}

	return client
}
