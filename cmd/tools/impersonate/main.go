package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/woningfinder/woningfinder/internal/auth"
	"github.com/woningfinder/woningfinder/internal/customer"
	"github.com/woningfinder/woningfinder/pkg/config"
)

// init is invoked before main()
func init() {
	// loads values from .env into the system
	// fallback to system env if unexisting
	// if not defined on system, panics
	if err := godotenv.Load("../../../.env"); err != nil {
		_ = config.MustGetString("APP_NAME")
	}
}

func main() {
	if len(os.Args) != 3 {
		log.Fatal("usage impersonate userID userEmail")
	}

	userID, err := strconv.ParseUint(os.Args[1], 10, 64)
	if err != nil {
		log.Fatal("userID must be an interger")
	}
	userEmail := os.Args[2]

	_, token, _ := auth.CreateJWTUserToken(auth.CreateJWTAuthenticationToken(config.MustGetString("JWT_SECRET")), &customer.User{
		ID:    uint(userID),
		Email: userEmail,
	})

	fmt.Printf("JWT token for %s: %s\n", userEmail, token)
}
