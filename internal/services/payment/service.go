package payment

import (
	"github.com/woningfinder/woningfinder/internal/customer"
	"github.com/woningfinder/woningfinder/internal/database"
	notificationService "github.com/woningfinder/woningfinder/internal/services/notification"
	userService "github.com/woningfinder/woningfinder/internal/services/user"
	"github.com/woningfinder/woningfinder/pkg/logging"
)

// Service permits to handle the management of the payments
type Service interface {
	ProcessPayment(email string, plan customer.Plan) error
}

type service struct {
	logger              *logging.Logger
	redisClient         database.RedisClient
	userService         userService.Service
	notificationService notificationService.Service
}

// NewService instantiate the payment service
func NewService(logger *logging.Logger, redisClient database.RedisClient, userService userService.Service, notificationService notificationService.Service) Service {
	return &service{
		logger:              logger,
		redisClient:         redisClient,
		userService:         userService,
		notificationService: notificationService,
	}
}
