package bootstrap

import (
	"net/http"
	"time"

	"github.com/woningfinder/woningfinder/internal/corporation/dewoonplaats"
	"github.com/woningfinder/woningfinder/pkg/networking"
	"github.com/woningfinder/woningfinder/pkg/networking/middleware"
)

// CreateDWClient creates a client for the De Woonplaats API
func CreateDWClient() dewoonplaats.Client {
	client := &http.Client{
		Timeout: 3 * time.Second,
	}
	defaultMiddleWare := []networking.ClientMiddleware{
		middleware.CreateHostMiddleware(dewoonplaats.Host),
		middleware.CreateDefaultHeadersMiddleware(map[string]string{"Content-Type": "application/json"}),
	}

	httpClient := networking.NewClient(client, defaultMiddleWare...)

	return dewoonplaats.NewClient(httpClient)
}
