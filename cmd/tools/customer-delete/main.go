package main

import (
	"github.com/woningfinder/woningfinder/pkg/logging"

	"github.com/woningfinder/woningfinder/internal/bootstrap"
	"github.com/woningfinder/woningfinder/internal/corporation"
	"github.com/woningfinder/woningfinder/internal/user"
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
}

func main() {
	logger := logging.NewZapLogger()

	err := bootstrap.InitDB()
	if err != nil {
		logger.Sugar().Fatal(err)
	}

	clientProvider := bootstrap.CreateClientProvider(logger, nil)
	corporationService := corporation.NewService(logger, bootstrap.DB, nil)
	userService := user.NewService(logger, bootstrap.DB, bootstrap.RDB, config.MustGetString("AES_SECRET"), clientProvider, corporationService)

	// get user
	u, err := userService.GetUser("PLACEHOLDER_EMAIL_TO_DELETE")
	if err != nil {
		logger.Sugar().Fatal(err)
	}

	// delete user
	if err := userService.DeleteUser(u); err != nil {
		logger.Sugar().Fatal(err)
	}

	logger.Sugar().Infof("customer %s successfully deleted 😢\n", u.Name)
}
