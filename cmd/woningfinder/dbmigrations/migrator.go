package dbmigrations

import (
	"fmt"

	"github.com/go-pg/migrations/v8"
	"github.com/julienrbrt/woningfinder/internal/bootstrap"
	bootstrapCorporation "github.com/julienrbrt/woningfinder/internal/bootstrap/corporation"
	"github.com/julienrbrt/woningfinder/internal/corporation/city"
	"github.com/julienrbrt/woningfinder/internal/corporation/connector"
	"github.com/julienrbrt/woningfinder/internal/database"
	"github.com/julienrbrt/woningfinder/internal/services/corporation"
	"github.com/julienrbrt/woningfinder/pkg/config"
	"github.com/julienrbrt/woningfinder/pkg/logging"
	"github.com/julienrbrt/woningfinder/pkg/mapbox"
	"go.uber.org/zap"
)

var (
	migrationsDbClient database.DBClient
	corporationService corporation.Service
	connectorProvider  connector.ConnectorProvider
)

func init() {
	logger := logging.NewZapLogger(config.GetBoolOrDefault("APP_DEBUG", false), config.GetStringOrDefault("SENTRY_DSN", ""))
	connectorProvider = bootstrapCorporation.CreateConnectorProvider(logger, mapbox.NewClientMock(nil, ""))
	migrationsDbClient = bootstrap.CreateDBClient(logger)
	corporationService = corporation.NewService(logger, migrationsDbClient, city.NewSuggester(connectorProvider.GetCities()))

}

// RunMigrations runs the database migrations.
func RunMigrations(logger *logging.Logger, dbClient database.DBClient, init bool) error {
	if init {
		// initialize database
		if _, _, err := migrations.Run(dbClient.Conn(), "init"); err != nil {
			return fmt.Errorf("error while initializing database %w", err)
		}
	}

	// run migrations
	oldVersion, newVersion, err := migrations.Run(dbClient.Conn())
	if err != nil {
		return fmt.Errorf("error while migrating database %w", err)
	}

	if newVersion != oldVersion {
		logger.Info("database schema migrated", zap.Int64("old_version", oldVersion), zap.Int64("new_version", newVersion))
	} else {
		logger.Info("database schema not updated", zap.Int64("version", oldVersion))
	}

	// close migrations db client
	if err := migrationsDbClient.Conn().Close(); err != nil {
		logger.Error("error closing migrations db client", zap.Error(err))
	}

	return nil
}
