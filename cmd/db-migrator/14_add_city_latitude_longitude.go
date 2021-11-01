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
		fmt.Println("adding longitude and latitude columbs column...")
		_, err := db.Exec(`ALTER TABLE cities ADD IF NOT EXISTS latitude float8`)
		if err != nil {
			return err
		}

		_, err = db.Exec(`ALTER TABLE cities ADD IF NOT EXISTS longitude float8`)
		if err != nil {
			return err
		}

		return nil
	})
}
