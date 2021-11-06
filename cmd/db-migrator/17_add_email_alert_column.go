package main

import (
	"github.com/go-pg/migrations/v8"
)

func init() {
	migrations.MustRegisterTx(func(db migrations.DB) error {
		logger.Info("adding user has_alerts_enabled column...")
		_, err := db.Exec(`ALTER TABLE users ADD IF NOT EXISTS has_alerts_enabled boolean`)
		return err
	})
}
