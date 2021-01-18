package bootstrap

import (
	"os"

	"github.com/go-pg/pg/v10/orm"
	"github.com/woningfinder/woningfinder/internal/database"
	"github.com/woningfinder/woningfinder/internal/domain/entity"
	"github.com/woningfinder/woningfinder/pkg/config"
	"github.com/woningfinder/woningfinder/pkg/logging"
)

// CreateDBClient creates a client for PostgreSQL and migrates the database upon creation
func CreateDBClient(logger *logging.Logger) database.DBClient {
	client, err := database.NewDBClient(logger, config.GetBoolOrDefault("APP_DEBUG", false), config.MustGetString("PSQL_HOST"), os.Getenv("PSQL_PORT"), os.Getenv("PSQL_NAME"), os.Getenv("PSQL_USERNAME"), os.Getenv("PSQL_PASSWORD"))
	if err != nil {
		logger.Sugar().Fatal(err)
	}

	return client
}

// RegisterModel many to many model so ORM can better recognize m2m relation.
// This should be done before dependant models are used.
func RegisterModel() {
	orm.RegisterTable((*entity.CorporationCity)(nil))
	orm.RegisterTable((*entity.CorporationSelectionMethod)(nil))
	orm.RegisterTable((*entity.HousingPreferencesHousingType)(nil))
	orm.RegisterTable((*entity.HousingPreferencesCity)(nil))
	orm.RegisterTable((*entity.HousingPreferencesCityDistrict)(nil))
}
