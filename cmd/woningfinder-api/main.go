package main

import (
	"net/http"

	"github.com/joho/godotenv"
	"github.com/julienrbrt/woningfinder/internal/auth"
	"github.com/julienrbrt/woningfinder/internal/bootstrap"
	bootstrapCorporation "github.com/julienrbrt/woningfinder/internal/bootstrap/corporation"
	"github.com/julienrbrt/woningfinder/internal/corporation/city"
	"github.com/julienrbrt/woningfinder/internal/handler"
	"github.com/julienrbrt/woningfinder/internal/services/corporation"
	emailService "github.com/julienrbrt/woningfinder/internal/services/email"
	userService "github.com/julienrbrt/woningfinder/internal/services/user"
	"github.com/julienrbrt/woningfinder/pkg/config"
	"github.com/julienrbrt/woningfinder/pkg/logging"
	"go.uber.org/zap"
)

// init is invoked before main()
func init() {
	// loads values from .env into the system
	// fallback to system env if unexisting
	// if not defined on system, panics
	if err := godotenv.Load(); err != nil {
		_ = config.MustGetString("APP_NAME")
	}
}

func main() {
	logger := logging.NewZapLogger(config.GetBoolOrDefault("APP_DEBUG", false), config.GetStringOrDefault("SENTRY_DSN", ""))
	dbClient := bootstrap.CreateDBClient(logger)
	jwtAuth := auth.CreateJWTAuthenticationToken(config.MustGetString("JWT_SECRET"))
	emailClient := bootstrap.CreateEmailClient()
	imgClient := bootstrap.CreateImgDownloader(logger)

	connectorProvider := bootstrapCorporation.CreateConnectorProvider(logger, nil) // mapboxClient not required in the api
	corporationService := corporation.NewService(logger, dbClient, city.NewSuggester(connectorProvider.GetCities()))
	userService := userService.NewService(logger, dbClient, config.MustGetString("AES_SECRET"), connectorProvider, corporationService)
	emailService := emailService.NewService(logger, emailClient, jwtAuth, imgClient)
	router := handler.NewHandler(logger, jwtAuth, corporationService, userService, emailService, imgClient)

	if err := http.ListenAndServe(":"+config.MustGetString("APP_PORT"), router); err != nil {
		logger.Fatal("failed to start server", zap.Error(err))
	}
}
