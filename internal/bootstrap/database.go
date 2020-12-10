package bootstrap

import (
	"fmt"
	"log"
	"os"

	"github.com/woningfinder/woningfinder/internal/corporation"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

//DB stores the database connection
var DB *gorm.DB

// InitDB create a connection to WoningFinder database
func InitDB() error {
	dbHost := os.Getenv("PSQL_HOST")
	dbUser := os.Getenv("PSQL_USERNAME")
	dbPassword := os.Getenv("PSQL_PASSWORD")
	dbName := os.Getenv("PSQL_NAME")
	dbPort := os.Getenv("PSQL_PORT")

	// build connection string
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=require TimeZone=Europe/Amsterdam", dbHost, dbUser, dbPassword, dbName, dbPort)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	DB = db
	if DB != nil {
		log.Printf("connected to database with host: %s\n", dbHost)
	}

	// Migrate the schema
	DB.Debug().AutoMigrate(&corporation.Corporation{}, &corporation.SelectionMethod{}, &corporation.City{})

	return nil
}
