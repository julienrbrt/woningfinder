package main

import (
	"github.com/go-pg/migrations/v8"
	"github.com/woningfinder/woningfinder/internal/corporation/connector/woningnet"
)

func init() {
	migrations.MustRegisterTx(func(db migrations.DB) error {
		return corporationService.CreateOrUpdateCorporation(woningnet.HengeloBorneInfo)
	})
}
