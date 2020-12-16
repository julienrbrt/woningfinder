package main

import (
	"log"
	"net/http"

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

	// app port
	port := env.MustGetString("APP_PORT")
	log.Println("listening on port", port)

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
