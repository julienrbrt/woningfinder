package main

import (
	"github.com/go-pg/migrations/v8"
)

func init() {

	migrations.MustRegisterTx(func(db migrations.DB) error {
		logger.Info("adding purchased_at column...")
		_, err := db.Exec(`ALTER TABLE user_plans ADD IF NOT EXISTS purchased_at timestamptz`)
		if err != nil {
			return err
		}

		_, err = db.Exec(`ALTER TABLE user_plans ADD IF NOT EXISTS free_trial_started_at timestamptz`)
		if err != nil {
			return err
		}

		return nil
	})
}
