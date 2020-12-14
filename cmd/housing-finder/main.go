package main

import (
	"log"
	"sync"

	"github.com/woningfinder/woningfinder/internal/bootstrap"
	"github.com/woningfinder/woningfinder/internal/corporation"
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

	err = bootstrap.InitRedis()
	if err != nil {
		log.Fatal(err)
	}

	clientProvider := bootstrap.CreateClientProvider()
	corporationService := corporation.NewService(bootstrap.DB, bootstrap.RDB)
	if err != nil {
		log.Fatal(err)
	}

	corporationList := clientProvider.List()
	wg := sync.WaitGroup{}
	for _, c := range *corporationList {
		client, err := clientProvider.Get(c)
		if err != nil {
			log.Println(err)
			continue
		}
		wg.Add(1)
		go func(corporation corporation.Corporation, client corporation.Client) {
			defer wg.Done()

			if err := corporationService.PublishOffers(client, corporation); err != nil {
				log.Println(err)
			}
		}(c, client)
	}
	wg.Wait()
}
