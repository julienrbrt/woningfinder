package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"
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

	// app port
	port := config.MustGetString("APP_PORT")
	log.Println("listening on port", port)

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
