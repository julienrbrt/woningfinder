package main

import (
	"log"

	"github.com/woningfinder/woningfinder/internal/corporation"
	"github.com/woningfinder/woningfinder/internal/corporation/onshuis"
	"gorm.io/gorm/clause"

	"github.com/woningfinder/woningfinder/internal/bootstrap"
	"github.com/woningfinder/woningfinder/internal/corporation/dewoonplaats"
	"github.com/woningfinder/woningfinder/pkg/env"

	"github.com/joho/godotenv"
)

// init is invoked before main()
func init() {
	// loads values from .env into the system
	// fallback to system env if unexisting
	// if not defined on system, panics
	if err := godotenv.Load("../../.env"); err != nil {
		_ = env.MustGetString("APP_NAME")
	}
}

func main() {
	err := bootstrap.InitDB()
	if err != nil {
		log.Fatal(err)
	}

	corporations := []corporation.Corporation{
		dewoonplaats.Info,
		onshuis.Info,
	}

	for _, c := range corporations {
		// creates the corporation - on data changes update it
		result := bootstrap.DB.Clauses(clause.OnConflict{
			UpdateAll: true,
		}).Create(&c)
		if result.Error != nil {
			log.Fatal(err)
		}
	}
}
