package main

import (
	"github.com/go-pg/migrations/v8"
)

func init() {
	migrations.MustRegisterTx(func(db migrations.DB) error {
		logger.Info("adding longitude and latitude columns...")
		_, err := db.Exec(`ALTER TABLE cities ADD IF NOT EXISTS latitude float8`)
		if err != nil {
			return err
		}

		_, err = db.Exec(`ALTER TABLE cities ADD IF NOT EXISTS longitude float8`)
		return err
	})
}
