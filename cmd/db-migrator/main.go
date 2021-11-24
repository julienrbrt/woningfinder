package main

import (
	"os"

	"github.com/go-pg/migrations/v8"
	"github.com/joho/godotenv"
	"github.com/julienrbrt/woningfinder/internal/bootstrap"
	bootstrapCorporation "github.com/julienrbrt/woningfinder/internal/bootstrap/corporation"
	"github.com/julienrbrt/woningfinder/internal/corporation/city"
	"github.com/julienrbrt/woningfinder/internal/corporation/connector"
	"github.com/julienrbrt/woningfinder/internal/database"
	"github.com/julienrbrt/woningfinder/internal/services/corporation"
	"github.com/julienrbrt/woningfinder/pkg/config"
	"github.com/julienrbrt/woningfinder/pkg/logging"
	"go.uber.org/zap"
)

var logger *logging.Logger
var dbClient database.DBClient
var corporationService corporation.Service
var connectorProvider connector.ConnectorProvider

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
	connectorProvider = bootstrapCorporation.CreateConnectorProvider(logger, nil)
	corporationService = corporation.NewService(logger, dbClient, city.NewSuggester(connectorProvider.GetCities()))

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
