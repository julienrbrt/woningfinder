package main

import (
	"github.com/go-pg/migrations/v8"
	"github.com/joho/godotenv"
	"github.com/woningfinder/woningfinder/internal/corporation"
	"github.com/woningfinder/woningfinder/internal/corporation/connector/ikwilhuren"
	"github.com/woningfinder/woningfinder/internal/corporation/connector/woonburo"
	"github.com/woningfinder/woningfinder/pkg/config"
)

func init() {
	// loads values from .env into the system
	// fallback to system env if unexisting
	// if not defined on system, panics
	if err := godotenv.Load("../../.env"); err != nil {
		_ = config.MustGetString("APP_NAME")
	}

	migrations.MustRegisterTx(func(db migrations.DB) error {
		for _, corp := range []corporation.Corporation{woonburo.AlmeloInfo, ikwilhuren.Info} {
			if err := corporationService.CreateOrUpdateCorporation(corp); err != nil {
				return err
			}
		}

		return nil
	})
}
