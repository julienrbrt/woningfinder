package main

import (
	"os"
	"strings"

	"github.com/woningfinder/woningfinder/internal/services/corporation"
	"github.com/woningfinder/woningfinder/internal/services/user"

	"github.com/woningfinder/woningfinder/pkg/logging"

	"github.com/woningfinder/woningfinder/internal/bootstrap"
	"github.com/woningfinder/woningfinder/pkg/config"

	"github.com/joho/godotenv"
)

// init is invoked before main()
func init() {
	// loads values from .env into the system
	// fallback to system env if unexisting
	// if not defined on system, panics
	if err := godotenv.Load("../../../.env"); err != nil {
		_ = config.MustGetString("APP_NAME")
	}

	// register m2m models with go-pg
	bootstrap.RegisterModel()
}

func main() {
	logger := logging.NewZapLogger(config.GetBoolOrDefault("APP_DEBUG", false), config.MustGetString("SENTRY_DSN"))

	// read email to delete from arguments
	if len(os.Args) != 2 {
		logger.Sugar().Fatal("customer-delete must have an user email as (only) argument\n")
	}
	email := os.Args[0]
	if email == "" || !strings.Contains(email, "@") {
		logger.Sugar().Fatal("incorrect argument for user to delete, have %s, expect a correct email", email)
	}

	dbClient := bootstrap.CreateDBClient(logger)
	clientProvider := bootstrap.CreateClientProvider(logger, nil)
	corporationService := corporation.NewService(logger, dbClient, nil)
	userService := user.NewService(logger, dbClient, nil, config.MustGetString("AES_SECRET"), clientProvider, corporationService)

	// get user
	u, err := userService.GetUser(email)
	if err != nil {
		logger.Sugar().Fatal(err)
	}

	// delete user
	if err := userService.DeleteUser(u); err != nil {
		logger.Sugar().Fatal(err)
	}

	logger.Sugar().Infof("customer %s successfully deleted 😢\n", u.Name)
}
