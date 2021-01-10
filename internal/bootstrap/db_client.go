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
	// DB.Debug().AutoMigrate(...) for extensive log
	client.Conn().AutoMigrate(
		&entity.Corporation{},
		&entity.SelectionMethod{},
		&entity.HousingType{},
		&entity.City{},
		&entity.CityDistrict{},
		&entity.User{},
		&entity.HousingPreferences{},
		&entity.CorporationCredentials{},
	)

	return client
}
