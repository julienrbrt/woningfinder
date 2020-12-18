package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/woningfinder/woningfinder/internal/bootstrap"
	"github.com/woningfinder/woningfinder/internal/corporation"
	"github.com/woningfinder/woningfinder/internal/user"
	"github.com/woningfinder/woningfinder/pkg/config"
)

// init is invoked before main()
func init() {
	// loads values from .env into the system
	// fallback to system env if unexisting
	// if not defined on system, panics
	if err := godotenv.Load("../../.env"); err != nil {
		_ = config.MustGetString("APP_NAME")
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
	userService := user.NewService(bootstrap.DB, bootstrap.RDB, config.MustGetString("AES_SECRET"), clientProvider, corporationService)

	offerList := make(chan corporation.OfferList)
	// subscribe to pub/sub messages inside a new goroutine
	go corporationService.SubscribeOffers(offerList)

	for o := range offerList {
		if err := userService.MatchOffer(o); err != nil {
			log.Printf("error while maching offers for corporation %s: %v\n", o.Corporation.Name, err)
		}
	}

}
