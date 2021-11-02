package main

import (
	"github.com/go-pg/migrations/v8"
	"github.com/joho/godotenv"
	"github.com/woningfinder/woningfinder/internal/bootstrap"
	"github.com/woningfinder/woningfinder/internal/corporation"
	"github.com/woningfinder/woningfinder/internal/corporation/connector/woningnet"
	corporationService "github.com/woningfinder/woningfinder/internal/services/corporation"
	"github.com/woningfinder/woningfinder/pkg/config"
	"github.com/woningfinder/woningfinder/pkg/logging"
)

func init() {
	// loads values from .env into the system
	// fallback to system env if unexisting
	// if not defined on system, panics
	if err := godotenv.Load("../../.env"); err != nil {
		_ = config.MustGetString("APP_NAME")
	}

	logger := logging.NewZapLoggerWithoutSentry()
	dbClient := bootstrap.CreateDBClient(logger)
	corporationService := corporationService.NewService(logger, dbClient)

	migrations.MustRegisterTx(func(db migrations.DB) error {
		for _, corp := range []corporation.Corporation{
			woningnet.AlmereInfo,
			woningnet.WoonkeusInfo,
			woningnet.EemvalleiInfo,
			woningnet.WoonserviceInfo,
			woningnet.MercatusInfo,
			woningnet.MiddenHollandInfo,
			woningnet.BovenGroningenInfo,
			woningnet.GooiVechtstreekInfo,
			woningnet.GroningenInfo,
			woningnet.HuiswaartsInfo,
			woningnet.WoongaardInfo,
		} {
			if err := corporationService.CreateOrUpdateCorporation(corp); err != nil {
				return err
			}
		}

		return nil
	})
}
