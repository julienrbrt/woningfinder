package bootstrap

import (
	"github.com/julienrbrt/woningfinder/internal/database"
	"github.com/julienrbrt/woningfinder/pkg/config"
	"github.com/julienrbrt/woningfinder/pkg/logging"
	"go.uber.org/zap"
)

// CreateDBClient creates a client for PostgreSQL and migrates the database upon creation
func CreateDBClient(logger *logging.Logger) database.DBClient {
	client, err := database.NewDBClient(logger, config.GetBoolOrDefault("APP_DEBUG", false), config.MustGetString("DATABASE_URL"))
	if err != nil {
		logger.Fatal("error creating database client", zap.Error(err))
	}

	return client
}
