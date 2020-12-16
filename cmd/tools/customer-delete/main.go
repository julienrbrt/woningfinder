package main

import (
	"log"
	"os"

	"github.com/woningfinder/woningfinder/internal/bootstrap"
	"github.com/woningfinder/woningfinder/internal/corporation"
	"github.com/woningfinder/woningfinder/internal/user"
	"github.com/woningfinder/woningfinder/pkg/env"

	"github.com/joho/godotenv"
)

// init is invoked before main()
func init() {
	// loads values from .env into the system
	// fallback to system env if unexisting
	// if not defined on system, panics
	if err := godotenv.Load("../../../.env"); err != nil {
		_ = env.MustGetString("APP_NAME")
	}
}

func main() {
	err := bootstrap.InitDB()
	if err != nil {
		log.Fatal(err)
	}

	clientProvider := bootstrap.CreateClientProvider()
	corporationService := corporation.NewService(bootstrap.DB, nil)
	userService := user.NewService(bootstrap.DB, bootstrap.RDB, os.Getenv("AES_SECRET"), clientProvider, corporationService)

	// get user
	u, err := userService.GetUser("PLACEHOLDER_EMAIL_TO_DELETE")
	if err != nil {
		log.Fatal(err)
	}

	// delete user
	if err := userService.DeleteUser(u); err != nil {
		log.Fatal(err)
	}

	log.Printf("customer %s successfully deleted 😢\n", u.Name)
}
