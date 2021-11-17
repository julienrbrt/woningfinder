package main

import (
	"github.com/go-pg/migrations/v8"
	bootstrapCorporation "github.com/woningfinder/woningfinder/internal/bootstrap/corporation"
)

func init() {
	corporations := bootstrapCorporation.CreateClientProvider(logger, nil).GetCorporations()
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
