package main

import (
	"github.com/go-pg/migrations/v8"
	"github.com/joho/godotenv"
	"github.com/woningfinder/woningfinder/internal/bootstrap"
	bootstrapCorporation "github.com/woningfinder/woningfinder/internal/bootstrap/corporation"
	"github.com/woningfinder/woningfinder/internal/services/corporation"
	"github.com/woningfinder/woningfinder/pkg/config"
	"github.com/woningfinder/woningfinder/pkg/logging"
)

func init() {
	// loads values from .env into the system
	// fallback to system env if unexisting
	// if not defined on system, panics
	if err := godotenv.Load("../../../.env"); err != nil {
		_ = config.MustGetString("APP_NAME")
	}

	logger := logging.NewZapLoggerWithoutSentry()
	corporations := bootstrapCorporation.CreateClientProvider(logger, nil).List()
	dbClient := bootstrap.CreateDBClient(logger)
	corporationService := corporation.NewService(logger, dbClient)

	migrations.MustRegisterTx(func(db migrations.DB) error {
		// add corporations
		for _, corp := range corporations {
			if err := corporationService.CreateOrUpdateCorporation(corp); err != nil {
				return err
			}
		}

		return nil
	})
}
