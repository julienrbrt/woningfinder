package bootstrap

import (
	"os"

	"github.com/woningfinder/woningfinder/internal/domain/entity"

	"github.com/woningfinder/woningfinder/internal/database"
	"github.com/woningfinder/woningfinder/pkg/config"
	"go.uber.org/zap"
)

// CreateDBClient creates a client for PostgreSQL and migrates the database upon creation
func CreateDBClient(logger *zap.Logger) database.DBClient {
	client, err := database.NewDBClient(logger, config.MustGetString("PSQL_HOST"), os.Getenv("PSQL_PORT"), os.Getenv("PSQL_NAME"), os.Getenv("PSQL_USERNAME"), os.Getenv("PSQL_PASSWORD"))
	if err != nil {
		logger.Sugar().Fatal(err)
	}

	// Migrate the schema
	db := client.Conn()
	if config.GetBoolOrDefault("APP_DEBUG", false) {
		db = db.Debug()
	}

	if err = db.AutoMigrate(
		&entity.Corporation{},
		&entity.SelectionMethod{},
		&entity.HousingType{},
		&entity.City{},
		&entity.CityDistrict{},
		&entity.User{},
		&entity.Plan{},
		&entity.HousingPreferences{},
		&entity.HousingPreferencesMatch{},
		&entity.CorporationCredentials{},
	); err != nil {
		logger.Sugar().Fatal(err)
	}

	return client
}
