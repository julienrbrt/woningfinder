package main

import (
	"net/http"

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
}

func main() {
	logger := logging.NewZapLoggerWithSentry(config.MustGetString("SENTRY_DSN"))

	// app port
	port := config.MustGetString("APP_PORT")
	logger.Sugar().Infof("listening on port %d", port)

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		logger.Sugar().Fatalf("failed to start server: %w", err)
	}
}
