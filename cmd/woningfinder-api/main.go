package main

import (
	"net/http"

	"github.com/woningfinder/woningfinder/internal/services/user"

	"github.com/woningfinder/woningfinder/internal/bootstrap"
	"github.com/woningfinder/woningfinder/internal/handler"
	"github.com/woningfinder/woningfinder/internal/services/corporation"

	"github.com/woningfinder/woningfinder/pkg/logging"

	"github.com/joho/godotenv"
	"github.com/woningfinder/woningfinder/pkg/config"
)

// init is invoked before main()
func init() {
	// loads values from .env into the system
	// fallback to system env if unexisting
	// if not defined on system, panics
	if err := godotenv.Load("../../.env"); err != nil {
		_ = config.MustGetString("APP_NAME")
	}

	// register m2m models with go-pg
	bootstrap.RegisterModel()
}

func main() {
	logger := logging.NewZapLogger(config.GetBoolOrDefault("APP_DEBUG", false), config.MustGetString("SENTRY_DSN"))
	dbClient := bootstrap.CreateDBClient(logger)
	redisClient := bootstrap.CreateRedisClient(logger)
	corporationService := corporation.NewService(logger, dbClient)
	clientProvider := bootstrap.CreateClientProvider(logger, nil) // mapboxClient not required in the api
	userService := user.NewService(logger, dbClient, redisClient, config.MustGetString("AES_SECRET"), clientProvider, corporationService)
	router := handler.NewHandler(logger, corporationService, userService)

	if err := http.ListenAndServe(":"+config.MustGetString("APP_PORT"), router); err != nil {
		logger.Sugar().Fatalf("failed to start server: %w", err)
	}
}
