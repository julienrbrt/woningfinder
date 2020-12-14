package main

import (
	"log"

	"github.com/woningfinder/woningfinder/internal/corporation"

	"github.com/woningfinder/woningfinder/internal/bootstrap"
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

var housingTypes = []corporation.HousingType{
	{
		Type: corporation.House,
	},
	{
		Type: corporation.Appartement,
	},
	{
		Type: corporation.Parking,
	},
	{
		Type: corporation.Undefined,
	},
}

func main() {
	err := bootstrap.InitDB()
	if err != nil {
		log.Fatal(err)
	}

	clientProvider := bootstrap.CreateClientProvider()
	corporationService := corporation.NewService(bootstrap.DB, nil)

	if _, err := corporationService.CreateOrUpdate(clientProvider.List()); err != nil {
		log.Fatal(err)
	}

	if _, err := corporationService.CreateHousingType(&housingTypes); err != nil {
		log.Fatal(err)
	}

	log.Println("successfully initialized database 🎉")
}
