package main

import (
	"github.com/go-pg/migrations/v8"
	"github.com/joho/godotenv"
	"github.com/julienrbrt/woningfinder/pkg/config"
)

func init() {
	// loads values from .env into the system
	// fallback to system env if unexisting
	// if not defined on system, panics
	if err := godotenv.Load("../../.env"); err != nil {
		_ = config.MustGetString("APP_NAME")
	}

	migrations.MustRegisterTx(func(db migrations.DB) error {
		_, err := db.Exec(`ALTER TABLE user_plans RENAME COLUMN purchased_at TO subscription_started_at`)
		if err != nil {
			return err
		}

		_, err = db.Exec(`ALTER TABLE user_plans RENAME COLUMN free_trial_started_at TO activated_at`)
		if err != nil {
			return err
		}

		_, err = db.Exec(`ALTER TABLE user_plans ADD IF NOT EXISTS last_payment_succeeded boolean`)
		if err != nil {
			return err
		}

		_, err = db.Exec(`ALTER TABLE user_plans ADD IF NOT EXISTS stripe_customer_id text`)
		if err != nil {
			return err
		}

		return nil
	})
}
