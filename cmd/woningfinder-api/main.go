package main

import (
	"net/http"

	"github.com/joho/godotenv"
	"github.com/woningfinder/woningfinder/internal/auth"
	"github.com/woningfinder/woningfinder/internal/bootstrap"
	bootstrapCorporation "github.com/woningfinder/woningfinder/internal/bootstrap/corporation"
	"github.com/woningfinder/woningfinder/internal/handler"
	"github.com/woningfinder/woningfinder/internal/services/corporation"
	emailService "github.com/woningfinder/woningfinder/internal/services/email"
	userService "github.com/woningfinder/woningfinder/internal/services/user"
	"github.com/woningfinder/woningfinder/pkg/config"
	"github.com/woningfinder/woningfinder/pkg/logging"
	"go.uber.org/zap"
)

// init is invoked before main()
func init() {
	// loads values from .env into the system
	// fallback to system env if unexisting
	// if not defined on system, panics
	if err := godotenv.Load("../../.env"); err != nil {
		_ = config.MustGetString("APP_NAME")
	}
}

func main() {
	logger := logging.NewZapLogger(config.GetBoolOrDefault("APP_DEBUG", false), config.MustGetString("SENTRY_DSN"))
	dbClient := bootstrap.CreateDBClient(logger)
	jwtAuth := auth.CreateJWTAuthenticationToken(config.MustGetString("JWT_SECRET"))
	emailClient := bootstrap.CreateEmailClient()
	stripeClient := bootstrap.CreateSripeClient(logger) // init stripe library

	corporationService := corporation.NewService(logger, dbClient)
	clientProvider := bootstrapCorporation.CreateClientProvider(logger, nil) // mapboxClient not required in the api
	userService := userService.NewService(logger, dbClient, config.MustGetString("AES_SECRET"), clientProvider, corporationService)
	emailService := emailService.NewService(logger, emailClient, jwtAuth)
	router := handler.NewHandler(logger, jwtAuth, corporationService, userService, emailService, stripeClient)

	if err := http.ListenAndServe(":"+config.MustGetString("APP_PORT"), router); err != nil {
		logger.Fatal("failed to start server", zap.Error(err))
	}
}
