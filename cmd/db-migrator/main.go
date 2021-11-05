package main

import (
	"os"

	"github.com/go-pg/migrations/v8"
	"github.com/joho/godotenv"
	"github.com/woningfinder/woningfinder/internal/bootstrap"
	"github.com/woningfinder/woningfinder/internal/database"
	"github.com/woningfinder/woningfinder/internal/services/corporation"
	"github.com/woningfinder/woningfinder/pkg/config"
	"github.com/woningfinder/woningfinder/pkg/logging"
	"go.uber.org/zap"
)

var logger *logging.Logger
var dbClient database.DBClient
var corporationService corporation.Service

// init is invoked before main()
func init() {
	// loads values from .env into the system
	// fallback to system env if unexisting
	// if not defined on system, panics
	if err := godotenv.Load("../../.env"); err != nil {
		_ = config.MustGetString("APP_NAME")
	}

	logger = logging.NewZapLogger(config.GetBoolOrDefault("APP_DEBUG", false), config.MustGetString("SENTRY_DSN"))
	dbClient = bootstrap.CreateDBClient(logger)
	corporationService = corporation.NewService(logger, dbClient)
}

func main() {
	if len(os.Args) > 1 && os.Args[1] == "init" {
		// initialize database
		if _, _, err := migrations.Run(dbClient.Conn(), "init"); err != nil {
			logger.Fatal("error while initializing database", zap.Error(err))
		}
	}

	// run migrations
	oldVersion, newVersion, err := migrations.Run(dbClient.Conn())
	if err != nil {
		logger.Fatal("error while migrating database", zap.Error(err))
	}

	if newVersion != oldVersion {
		logger.Info("database schema migrated", zap.Int64("old_version", oldVersion), zap.Int64("new_version", newVersion))
	} else {
		logger.Info("database schema not updated", zap.Int64("version", oldVersion))
	}
}
