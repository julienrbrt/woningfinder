package bootstrap

import (
	"fmt"
	"log"
	"os"

	"github.com/woningfinder/woningfinder/internal/user"

	"github.com/woningfinder/woningfinder/internal/corporation"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

//DB stores the database connection
var DB *gorm.DB

// InitDB create a connection to WoningFinder PostgreSQL database
func InitDB() error {
	dbHost := os.Getenv("PSQL_HOST")
	dbPort := os.Getenv("PSQL_PORT")
	dbName := os.Getenv("PSQL_NAME")
	dbUser := os.Getenv("PSQL_USERNAME")
	dbPassword := os.Getenv("PSQL_PASSWORD")

	// build connection string
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=require TimeZone=Europe/Amsterdam", dbHost, dbUser, dbPassword, dbName, dbPort)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	DB = db
	if DB != nil {
		log.Println("successfully connected to database")
	}

	// Migrate the schema
	// DB.Debug().AutoMigrate(...) for extensive log
	DB.AutoMigrate(
		&corporation.Corporation{},
		&corporation.SelectionMethod{},
		&corporation.City{},
		&corporation.District{},
		&corporation.HousingType{},
		&user.User{},
		&user.HousingPreferences{},
		&user.CorporationCredentials{},
	)

	return nil
}
