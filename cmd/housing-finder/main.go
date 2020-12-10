package main

import (
	"log"

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

func main() {
	err := bootstrap.InitDB()
	if err != nil {
		log.Fatal(err)
	}

}
