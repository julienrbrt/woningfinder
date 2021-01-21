package payment

import (
	"github.com/woningfinder/woningfinder/internal/database"
	"github.com/woningfinder/woningfinder/internal/domain/entity"
	userService "github.com/woningfinder/woningfinder/internal/services/user"
	"github.com/woningfinder/woningfinder/pkg/logging"
)

// pubSubPayment defines on which channel the payment confirmation are sent via redis
const pubSubPayment = "stripe"

// Service permits to handle the persistence of an user
type Service interface {
	QueuePayment(payment *entity.PaymentData) error
	ProcessPayment(paymentCh chan<- entity.PaymentData) error
}

type service struct {
	logger      *logging.Logger
	redisClient database.RedisClient
	userService userService.Service
}

// NewService instantiate the payment service
func NewService(logger *logging.Logger, redisClient database.RedisClient, userService userService.Service) Service {
	return &service{
		logger:      logger,
		redisClient: redisClient,
		userService: userService,
	}
}
