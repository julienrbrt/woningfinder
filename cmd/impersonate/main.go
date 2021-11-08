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
	"go.uber.org/zap"
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

const (
	baseURL    = "https://woningfinder.nl"
	baseDEVURL = "http://localhost:3000"
)

func main() {
	logger := logging.NewZapLoggerWithoutSentry()

	if len(os.Args) < 3 {
		logger.Fatal("usage impersonate userID email")
	}

	baseURL := baseURL
	if len(os.Args) > 3 {
		baseURL = baseDEVURL
	}

	userID, err := strconv.ParseUint(os.Args[1], 10, 64)
	if err != nil {
		logger.Fatal("userID must be an interger")
	}

	email := os.Args[2]
	if !util.IsEmailValid(email) {
		logger.Fatal("email invalid", zap.String("got", email))
	}

	_, token, _ := auth.CreateJWTUserToken(auth.CreateJWTAuthenticationToken(config.MustGetString("JWT_SECRET")), &customer.User{
		ID:    uint(userID),
		Email: email,
	})

	fmt.Printf("Authenticated with %s: %s/mijn-zoekopdracht?jwt=%s\n", email, baseURL, token)
}
