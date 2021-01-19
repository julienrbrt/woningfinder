package main

import (
	"github.com/go-pg/migrations/v8"
	"github.com/joho/godotenv"
	"github.com/woningfinder/woningfinder/internal/bootstrap"
	"github.com/woningfinder/woningfinder/internal/domain/entity"
	"github.com/woningfinder/woningfinder/internal/services/corporation"
	"github.com/woningfinder/woningfinder/pkg/config"
	"github.com/woningfinder/woningfinder/pkg/logging"
)

var housingTypes = []entity.HousingType{
	{
		Type: entity.HousingTypeHouse,
	},
	{
		Type: entity.HousingTypeAppartement,
	},
	{
		Type: entity.HousingTypeUndefined,
	},
}

var selectionMethods = []entity.SelectionMethod{
	{
		Method: entity.SelectionFirstComeFirstServed,
	},
	{
		Method: entity.SelectionRandom,
	},
	{
		Method: entity.SelectionRegistrationDate,
	},
}

func init() {
	// loads values from .env into the system
	// fallback to system env if unexisting
	// if not defined on system, panics
	if err := godotenv.Load("../../../.env"); err != nil {
		_ = config.MustGetString("APP_NAME")
	}

	logger := logging.NewZapLoggerWithoutSentry()
	corporations := bootstrap.CreateClientProvider(logger, nil).List()
	dbClient := bootstrap.CreateDBClient(logger)
	corporationService := corporation.NewService(logger, dbClient, nil)

	migrations.MustRegisterTx(func(db migrations.DB) error {
		// add housing types
		if _, err := db.Model(&housingTypes).Insert(); err != nil {
			return err
		}

		// add selections methods
		if _, err := db.Model(&selectionMethods).Insert(); err != nil {
			return err
		}

		// add corporations
		for _, corp := range corporations {
			if err := corporationService.CreateOrUpdateCorporation(corp); err != nil {
				return err
			}
		}

		return nil
	})
}
