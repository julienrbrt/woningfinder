package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/woningfinder/woningfinder/internal/auth"
	"github.com/woningfinder/woningfinder/internal/customer"
	"github.com/woningfinder/woningfinder/pkg/config"
	"github.com/woningfinder/woningfinder/pkg/logging"
	"github.com/woningfinder/woningfinder/pkg/util"
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
	logger := logging.NewZapLogger(config.GetBoolOrDefault("APP_DEBUG", false), config.MustGetString("SENTRY_DSN"))

	if len(os.Args) != 3 {
		logger.Sugar().Fatal("usage impersonate userID email")
	}

	userID, err := strconv.ParseUint(os.Args[1], 10, 64)
	if err != nil {
		logger.Sugar().Fatal("userID must be an interger")
	}

	email := os.Args[2]
	if !util.IsEmailValid(email) {
		logger.Sugar().Fatalf("email %s invalid", email)
	}

	_, token, _ := auth.CreateJWTUserToken(auth.CreateJWTAuthenticationToken(config.MustGetString("JWT_SECRET")), &customer.User{
		ID:    uint(userID),
		Email: email,
	})

	fmt.Printf("Authenticated with %s: https://woningfinder.nl/mijn-zoekopdracht?jwt=%s\n", email, token)
}
