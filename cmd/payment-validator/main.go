package main

import (
	"github.com/joho/godotenv"
	"github.com/woningfinder/woningfinder/internal/auth"
	"github.com/woningfinder/woningfinder/internal/bootstrap"
	"github.com/woningfinder/woningfinder/internal/domain/entity"
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
	paymentService := paymentService.NewService(logger, redisClient, userService)
	notificationsService := notificationsService.NewService(logger, emailClient, jwtAuth)

	payments := make(chan entity.PaymentData)
	// subscribe to pub/sub messages inside a new go routine
	go func(payments chan entity.PaymentData) {
		if err := paymentService.ProcessPayment(payments); err != nil {
			logger.Sugar().Fatal(err)
		}
	}(payments)

	// set that user has paid
	for payment := range payments {
		user, err := userService.GetUser(&entity.User{Email: payment.UserEmail})
		if err != nil {
			logger.Sugar().Errorf("error while processing payment data: cannot get user: %w", err)
			continue
		}

		if err := userService.SetPaid(user, payment.Plan); err != nil {
			logger.Sugar().Errorf("error while processing payment data: %w", err)
			continue
		}

		// send confirmation email
		if err := notificationsService.SendWelcome(user); err != nil {
			logger.Sugar().Errorf("error while processing payment data: %w", err)
		}
	}
}
