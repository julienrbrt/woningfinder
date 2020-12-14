package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/woningfinder/woningfinder/internal/bootstrap"
	"github.com/woningfinder/woningfinder/internal/corporation"
	"github.com/woningfinder/woningfinder/internal/user"
	"github.com/woningfinder/woningfinder/pkg/env"
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

	err = bootstrap.InitRedis()
	if err != nil {
		log.Fatal(err)
	}

	clientProvider := bootstrap.CreateClientProvider()
	corporationService := corporation.NewService(bootstrap.DB, bootstrap.RDB)
	userService := user.NewService(bootstrap.DB, os.Getenv("AES_SECRET"), clientProvider, corporationService)

	offer := make(chan corporation.Offer)
	// subscribe to pub/sub messages inside a new goroutine
	go corporationService.SubscribeOffers(offer)

	for o := range offer {
		if err := userService.MatchOffer(o); err != nil {
			log.Printf("error while maching offer %s: %v\n", o.Housing.Address, err)
		}
	}

}
