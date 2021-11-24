package main

import (
	"github.com/go-pg/migrations/v8"
	"github.com/julienrbrt/woningfinder/internal/corporation/connector/itris"
)

func init() {
	migrations.MustRegisterTx(func(db migrations.DB) error {
		return corporationService.CreateOrUpdateCorporation(itris.MijandeWonenInfo)
	})
}
