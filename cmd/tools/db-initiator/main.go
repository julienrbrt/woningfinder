package main

import (
	"github.com/joho/godotenv"
	"github.com/woningfinder/woningfinder/internal/bootstrap"
	"github.com/woningfinder/woningfinder/internal/corporation"
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

var housingTypes = []corporation.HousingType{
	{
		Type: corporation.House,
	},
	{
		Type: corporation.Appartement,
	},
	{
		Type: corporation.Undefined,
	},
}

func main() {
	logger := logging.NewZapLogger()

	err := bootstrap.InitDB()
	if err != nil {
		logger.Sugar().Fatal(err)
	}

	clientProvider := bootstrap.CreateClientProvider(logger)
	corporationService := corporation.NewService(logger, bootstrap.DB, nil)

	if _, err := corporationService.CreateOrUpdate(clientProvider.List()); err != nil {
		logger.Sugar().Fatal(err)
	}

	if _, err := corporationService.CreateHousingType(&housingTypes); err != nil {
		logger.Sugar().Fatal(err)
	}

	logger.Sugar().Infof("successfully initialized database 🎉")
}
