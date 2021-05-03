package main

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/woningfinder/woningfinder/internal/bootstrap"
	bootstrapCorporation "github.com/woningfinder/woningfinder/internal/bootstrap/corporation"
	"github.com/woningfinder/woningfinder/internal/customer"
	"github.com/woningfinder/woningfinder/internal/services/corporation"
	userService "github.com/woningfinder/woningfinder/internal/services/user"
	"github.com/woningfinder/woningfinder/pkg/config"
	"github.com/woningfinder/woningfinder/pkg/logging"
	"github.com/woningfinder/woningfinder/pkg/util"
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
	logger := logging.NewZapLogger(config.GetBoolOrDefault("APP_DEBUG", false), config.MustGetString("SENTRY_DSN"))

	// read email to delete from arguments
	if len(os.Args) != 2 {
		logger.Sugar().Fatal("customer-delete must have an user email as (only) argument\n")
	}
	email := os.Args[1]
	if !util.IsEmailValid(email) {
		logger.Sugar().Fatal("incorrect argument for user to delete, have %s, expect a valid email", email)
	}

	dbClient := bootstrap.CreateDBClient(logger)
	clientProvider := bootstrapCorporation.CreateClientProvider(logger, nil)
	corporationService := corporation.NewService(logger, dbClient)
	userService := userService.NewService(logger, dbClient, nil, config.MustGetString("AES_SECRET"), clientProvider, corporationService)

	// get user
	user, err := userService.GetUser(&customer.User{Email: email})
	if err != nil {
		logger.Sugar().Fatal(err)
	}

	// delete user
	if err := userService.DeleteUser(user); err != nil {
		logger.Sugar().Fatal(err)
	}

	logger.Sugar().Infof("customer %s (%s) successfully deleted 😢\n", user.Name, user.Email)
}
