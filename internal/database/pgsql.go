package database

import (
	"fmt"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"moul.io/zapgorm2"
)

type DBClient interface {
	Conn() *gorm.DB
}

type dbClient struct {
	logger *zap.Logger
	db     *gorm.DB
}

// NewDBClient create a connection to postgresql database
func NewDBClient(logger *zap.Logger, host, port, name, user, password string) (DBClient, error) {
	gormLogger := zapgorm2.New(logger)
	gormLogger.SetAsDefault() // configure gorm to use this zapgorm.Logger for callbacks

	// build connection string
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=require TimeZone=Europe/Amsterdam", host, user, password, name, port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: gormLogger})
	if err != nil {
		return nil, err
	} else if db != nil {
		logger.Info("successfully connected to postgresql 🎉")
	}

	return &dbClient{
		logger: logger,
		db:     db,
	}, nil
}

func (db *dbClient) Conn() *gorm.DB {
	return db.db
}
