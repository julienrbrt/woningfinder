package main

import (
	"os"

	"github.com/woningfinder/woningfinder/internal/entity"
	"github.com/woningfinder/woningfinder/internal/services/corporation"
	"github.com/woningfinder/woningfinder/internal/services/user"
	"github.com/woningfinder/woningfinder/pkg/logging"
	"github.com/woningfinder/woningfinder/pkg/util"

	"github.com/joho/godotenv"
	"github.com/woningfinder/woningfinder/internal/bootstrap"
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
	logger := logging.NewZapLogger(config.GetBoolOrDefault("APP_DEBUG", false), config.MustGetString("SENTRY_DSN"))

	// read email to delete from arguments
	if len(os.Args) != 2 {
		logger.Sugar().Fatal("customer-delete must have an user email as (only) argument\n")
	}
	email := os.Args[0]
	if !util.IsEmailValid(email) {
		logger.Sugar().Fatal("incorrect argument for user to delete, have %s, expect a valid email", email)
	}

	dbClient := bootstrap.CreateDBClient(logger)
	clientProvider := bootstrap.CreateClientProvider(logger, nil)
	corporationService := corporation.NewService(logger, dbClient)
	userService := user.NewService(logger, dbClient, nil, config.MustGetString("AES_SECRET"), clientProvider, corporationService)

	// get user
	u, err := userService.GetUser(&entity.User{Email: email})
	if err != nil {
		logger.Sugar().Fatal(err)
	}

	// delete user
	if err := userService.DeleteUser(u); err != nil {
		logger.Sugar().Fatal(err)
	}

	logger.Sugar().Infof("customer %s successfully deleted 😢\n", u.Name)
}
