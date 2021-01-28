package main

import (
	"github.com/joho/godotenv"
	"github.com/woningfinder/woningfinder/internal/auth"
	"github.com/woningfinder/woningfinder/internal/bootstrap"
	notificationsService "github.com/woningfinder/woningfinder/internal/services/notifications"
	paymentService "github.com/woningfinder/woningfinder/internal/services/payment"
	userService "github.com/woningfinder/woningfinder/internal/services/user"
	"github.com/woningfinder/woningfinder/pkg/config"
	"github.com/woningfinder/woningfinder/pkg/logging"
)

// init is invoked before main()
func init() {
	// loads values from .env into the system
	// fallback to system env if unexisting
	// if not defined on system, panics
	if err := godotenv.Load("../../.env"); err != nil {
		_ = config.MustGetString("APP_NAME")
	}

	// register m2m models with go-pg
	bootstrap.RegisterModel()
}

func main() {
	logger := logging.NewZapLogger(config.GetBoolOrDefault("APP_DEBUG", false), config.MustGetString("SENTRY_DSN"))
	jwtAuth := auth.CreateJWTAuthenticationToken(config.MustGetString("JWT_SECRET"))

	dbClient := bootstrap.CreateDBClient(logger)
	redisClient := bootstrap.CreateRedisClient(logger)
	emailClient := bootstrap.CreateEmailClient()
	userService := userService.NewService(logger, dbClient, redisClient, config.MustGetString("AES_SECRET"), nil, nil) // no corporation relation stuff in payment-validator
	notificationsService := notificationsService.NewService(logger, emailClient, jwtAuth)
	paymentService := paymentService.NewService(logger, redisClient, userService, notificationsService)

	// process payment
	if err := paymentService.ProcessPayment(); err != nil {
		logger.Sugar().Fatal(err)
	}
}
