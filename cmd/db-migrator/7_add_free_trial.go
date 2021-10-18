package main

import (
	"fmt"

	"github.com/go-pg/migrations/v8"
	"github.com/joho/godotenv"
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
		fmt.Println("adding purchased_at column...")
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
