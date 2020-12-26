package bootstrap

import (
	"fmt"
	"os"

	"moul.io/zapgorm2"

	"github.com/woningfinder/woningfinder/pkg/logging"

	"github.com/woningfinder/woningfinder/internal/user"

	"github.com/woningfinder/woningfinder/internal/corporation"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

//DB stores the database connection
var DB *gorm.DB

// InitDB create a connection to WoningFinder PostgreSQL database
func InitDB() error {
	logger := logging.NewZapLogger()
	gormLogger := zapgorm2.New(logger)
	gormLogger.SetAsDefault() // configure gorm to use this zapgorm.Logger for callbacks

	dbHost := os.Getenv("PSQL_HOST")
	dbPort := os.Getenv("PSQL_PORT")
	dbName := os.Getenv("PSQL_NAME")
	dbUser := os.Getenv("PSQL_USERNAME")
	dbPassword := os.Getenv("PSQL_PASSWORD")

	// build connection string
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=require TimeZone=Europe/Amsterdam", dbHost, dbUser, dbPassword, dbName, dbPort)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: gormLogger})
	if err != nil {
		return err
	}

	DB = db
	if DB != nil {
		logger.Info("successfully connected to postgresql 🎉")
	}

	// Migrate the schema
	// DB.Debug().AutoMigrate(...) for extensive log
	DB.AutoMigrate(
		&corporation.Corporation{},
		&corporation.SelectionMethod{},
		&corporation.HousingType{},
		&corporation.City{},
		&corporation.CityDistrict{},
		&user.User{},
		&user.HousingPreferences{},
		&user.CorporationCredentials{},
	)

	return nil
}
