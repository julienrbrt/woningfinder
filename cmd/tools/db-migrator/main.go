package main

import (
	"os"

	"github.com/go-pg/migrations/v8"
	"github.com/joho/godotenv"
	"github.com/woningfinder/woningfinder/internal/bootstrap"
	"github.com/woningfinder/woningfinder/pkg/config"
	"github.com/woningfinder/woningfinder/pkg/logging"
)

// init is invoked before main()
func init() {
	// loads values from .env into the system
	// fallback to system env if unexisting
	// if not defined on system, panics
	if err := godotenv.Load("../../../.env"); err != nil {
		_ = config.MustGetString("APP_NAME")
	}
}

func main() {
	logger := logging.NewZapLogger(config.GetBoolOrDefault("APP_DEBUG", false), config.MustGetString("SENTRY_DSN"))
	dbClient := bootstrap.CreateDBClient(logger)

	if len(os.Args) > 1 && os.Args[1] == "init" {
		// initialize database
		if _, _, err := migrations.Run(dbClient.Conn(), "init"); err != nil {
			logger.Sugar().Fatalf("error while initializing database: %w", err)
		}
	}

	// run migrations
	oldVersion, newVersion, err := migrations.Run(dbClient.Conn())
	if err != nil {
		logger.Sugar().Fatalf("error while migrating database: %w", err)
	}

	if newVersion != oldVersion {
		logger.Sugar().Infof("database schema migrated from version %d to %d\n", oldVersion, newVersion)
	} else {
		logger.Sugar().Infof("database schema not updated: version stays %d\n", oldVersion)
	}
}
