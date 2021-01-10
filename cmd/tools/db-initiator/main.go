package main

import (
	"github.com/joho/godotenv"
	"github.com/woningfinder/woningfinder/internal/bootstrap"
	"github.com/woningfinder/woningfinder/internal/domain/entity"
	"github.com/woningfinder/woningfinder/internal/services/corporation"
	"github.com/woningfinder/woningfinder/pkg/config"
	"github.com/woningfinder/woningfinder/pkg/logging"
)

// init is invoked before main()
func init() {
	// loads values from .env into the system
	// fallback to system env if unexisting
	// if not defined on system, panics
	if err := godotenv.Load("../../../.env"); err != nil {
		_ = config.MustGetString("APP_NAME")
	}
}

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

func main() {
	logger := logging.NewZapLogger()

	dbClient := bootstrap.CreateDBClient(logger)
	clientProvider := bootstrap.CreateClientProvider(logger, nil)
	corporationService := corporation.NewService(logger, dbClient, nil)

	if _, err := corporationService.CreateOrUpdateCorporation(clientProvider.List()); err != nil {
		logger.Sugar().Fatal(err)
	}

	if _, err := corporationService.CreateHousingType(&housingTypes); err != nil {
		logger.Sugar().Fatal(err)
	}

	logger.Sugar().Info("successfully initialized database 🎉")
}
