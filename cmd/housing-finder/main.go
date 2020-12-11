package main

import (
	"fmt"
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

	corporationService := corporation.NewService(bootstrap.DB)
	corporations, err := corporationService.GetCorporations()
	if err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup
	ch := make(chan error, len(*corporations))
	for _, corporation := range *corporations {
		client, err := corporationService.GetClient(&corporation)
		if err != nil {
			log.Printf("cannot find client for corporation %s\n", corporation.Name)
			continue
		}

		wg.Add(1)
		go fetchAndSend(client, &wg, ch, corporation)
	}
	wg.Wait()
	fmt.Println(<-ch)
}

func fetchAndSend(client corporation.Client, wg *sync.WaitGroup, ch chan error, corporation corporation.Corporation) {
	defer wg.Done()

	offers, err := client.FetchOffer()
	if err != nil {
		ch <- fmt.Errorf("error while fetching offers for %s: %w", corporation.Name, err)
	}

	// technically send all fetched offers to redis
	fmt.Println(offers)
}
