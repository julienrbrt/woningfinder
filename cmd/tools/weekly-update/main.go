package main

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/woningfinder/woningfinder/internal/auth"
	"github.com/woningfinder/woningfinder/internal/bootstrap"
	bootstrapCorporation "github.com/woningfinder/woningfinder/internal/bootstrap/corporation"
	"github.com/woningfinder/woningfinder/internal/customer"
	corporationService "github.com/woningfinder/woningfinder/internal/services/corporation"
	emailService "github.com/woningfinder/woningfinder/internal/services/email"
	userService "github.com/woningfinder/woningfinder/internal/services/user"
	"github.com/woningfinder/woningfinder/pkg/config"
	"github.com/woningfinder/woningfinder/pkg/logging"
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

	if len(os.Args) != 2 {
		log.Fatal("usage weekly-update userEmail")
	}
	userEmail := os.Args[1]
	logger.Sugar().Infof("send weekly update to %s", userEmail)

	jwtAuth := auth.CreateJWTAuthenticationToken(config.MustGetString("JWT_SECRET"))

	dbClient := bootstrap.CreateDBClient(logger)
	mapboxClient := bootstrap.CreateMapboxClient()
	emailClient := bootstrap.CreateEmailClient()

	clientProvider := bootstrapCorporation.CreateClientProvider(logger, mapboxClient)
	corporationService := corporationService.NewService(logger, dbClient)
	userService := userService.NewService(logger, dbClient, config.MustGetString("AES_SECRET"), clientProvider, corporationService)
	emailService := emailService.NewService(logger, emailClient, jwtAuth)

	users, err := userService.GetWeeklyUpdateUsers()
	if err != nil {
		logger.Sugar().Errorf("error while sending weekly update: %w", err)
	}

	var user *customer.User
	for _, u := range users {
		if strings.EqualFold(u.Email, userEmail) {
			user = u
			break
		}
	}

	// user has no corporation credentials and no match didn't react for them be we cannot send weekly update
	if len(user.HousingPreferencesMatch) == 0 && len(user.CorporationCredentials) == 0 {
		if err := emailService.SendCorporationCredentialsMissing(user); err != nil {
			logger.Sugar().Errorf("error while sending weekly update (credentials missing): %w", err)
		}
	}

	if err := emailService.SendWeeklyUpdate(user, user.HousingPreferencesMatch); err != nil {
		logger.Sugar().Errorf("error while sending weekly update: %w", err)
	}
}
