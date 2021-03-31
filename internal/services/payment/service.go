package payment

import (
	"github.com/woningfinder/woningfinder/internal/database"
	"github.com/woningfinder/woningfinder/internal/domain/entity"
	notificationsService "github.com/woningfinder/woningfinder/internal/services/notifications"
	userService "github.com/woningfinder/woningfinder/internal/services/user"
	"github.com/woningfinder/woningfinder/pkg/logging"
)

// Service permits to handle the management of the payments
type Service interface {
	ProcessPayment(email string, plan entity.Plan) error
}

type service struct {
	logger               *logging.Logger
	redisClient          database.RedisClient
	userService          userService.Service
	notificationsService notificationsService.Service
}

// NewService instantiate the payment service
func NewService(logger *logging.Logger, redisClient database.RedisClient, userService userService.Service, notificationsService notificationsService.Service) Service {
	return &service{
		logger:               logger,
		redisClient:          redisClient,
		userService:          userService,
		notificationsService: notificationsService,
	}
}
