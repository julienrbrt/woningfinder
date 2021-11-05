package main

import (
	"github.com/go-pg/migrations/v8"
	"github.com/joho/godotenv"
	"github.com/woningfinder/woningfinder/internal/corporation/connector/zig"
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
		return corporationService.CreateOrUpdateCorporation(zig.WoonnetHaaglanden)
	})
}
