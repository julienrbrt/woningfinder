package main

import (
	"net/http"

	"github.com/joho/godotenv"
	"github.com/woningfinder/woningfinder/internal/auth"
	"github.com/woningfinder/woningfinder/internal/bootstrap"
	"github.com/woningfinder/woningfinder/internal/handler"
	"github.com/woningfinder/woningfinder/internal/services/corporation"
	notificationsService "github.com/woningfinder/woningfinder/internal/services/notifications"
	paymentService "github.com/woningfinder/woningfinder/internal/services/payment"
	userService "github.com/woningfinder/woningfinder/internal/services/user"
	"github.com/woningfinder/woningfinder/pkg/config"
	"github.com/woningfinder/woningfinder/pkg/logging"
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
	redisClient := bootstrap.CreateRedisClient(logger)
	emailClient := bootstrap.CreateEmailClient()
	jwtAuth := auth.CreateJWTAuthenticationToken(config.MustGetString("JWT_SECRET"))

	corporationService := corporation.NewService(logger, dbClient)
	clientProvider := bootstrap.CreateClientProvider(logger, nil) // mapboxClient not required in the api
	userService := userService.NewService(logger, dbClient, redisClient, config.MustGetString("AES_SECRET"), clientProvider, corporationService)
	bootstrap.CreateSripeClient(logger) // init stripe library
	notificationsService := notificationsService.NewService(logger, emailClient, jwtAuth)
	paymentService := paymentService.NewService(logger, redisClient, userService, notificationsService)
	router := handler.NewHandler(logger, corporationService, userService, paymentService, config.MustGetString("STRIPE_WEBHOOK_SIGNING_KEY"), jwtAuth, emailClient)

	if err := http.ListenAndServe(":"+config.MustGetString("APP_PORT"), router); err != nil {
		logger.Sugar().Fatalf("failed to start server: %w", err)
	}
}
