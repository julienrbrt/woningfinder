package main

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/woningfinder/woningfinder/internal/bootstrap"
	bootstrapCorporation "github.com/woningfinder/woningfinder/internal/bootstrap/corporation"
	"github.com/woningfinder/woningfinder/internal/services/corporation"
	userService "github.com/woningfinder/woningfinder/internal/services/user"
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

func main() {
	logger := logging.NewZapLoggerWithoutSentry()

	// read email to delete from arguments
	if len(os.Args) != 2 {
		logger.Fatal("usage: customer-delete email")
	}

	email := os.Args[1]
	if !util.IsEmailValid(email) {
		logger.Fatal("incorrect argument: expects a valid email")
	}

	dbClient := bootstrap.CreateDBClient(logger)
	clientProvider := bootstrapCorporation.CreateClientProvider(logger, nil)
	corporationService := corporation.NewService(logger, dbClient)
	userService := userService.NewService(logger, dbClient, config.MustGetString("AES_SECRET"), clientProvider, corporationService)

	// delete user
	if err := userService.DeleteUser(email); err != nil {
		logger.Fatal("error when deleting user", zap.Error(err))
	}

	logger.Info("customer successfully deleted 😢", zap.String("email", email))
}
