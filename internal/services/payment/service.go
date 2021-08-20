package payment

import (
	"github.com/woningfinder/woningfinder/internal/customer"
	"github.com/woningfinder/woningfinder/internal/database"
	emailService "github.com/woningfinder/woningfinder/internal/services/email"
	userService "github.com/woningfinder/woningfinder/internal/services/user"
	"github.com/woningfinder/woningfinder/pkg/logging"
)

// Service permits to handle the management of the payments
type Service interface {
	ProcessFreeTrial(email string, plan customer.Plan) error
	ProcessPayment(email string, plan customer.Plan) error
}

type service struct {
	logger       *logging.Logger
	dbClient     database.DBClient
	userService  userService.Service
	emailService emailService.Service
}

// NewService instantiate the payment service
func NewService(logger *logging.Logger, dbClient database.DBClient, userService userService.Service, emailService emailService.Service) Service {
	return &service{
		logger:       logger,
		dbClient:     dbClient,
		userService:  userService,
		emailService: emailService,
	}
}
