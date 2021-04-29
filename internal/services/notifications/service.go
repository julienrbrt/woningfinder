package notifications

import (
	"github.com/go-chi/jwtauth"
	"github.com/woningfinder/woningfinder/internal/customer"
	"github.com/woningfinder/woningfinder/pkg/email"
	"github.com/woningfinder/woningfinder/pkg/logging"
)

// Service permits to handle the management of the notifications
type Service interface {
	SendLogin(user *customer.User) error
	SendWelcome(user *customer.User) error
	SendPaymentReminder(user *customer.User) error
	SendWeeklyUpdate(user *customer.User, housingMatch []customer.HousingPreferencesMatch) error
	SendCorporationCredentialsError(user *customer.User, corporationName string) error
	SendBye(user *customer.User) error
}

type service struct {
	logger      *logging.Logger
	emailClient email.Client
	jwtAuth     *jwtauth.JWTAuth
}

// NewService instantiate the notification service
func NewService(logger *logging.Logger, emailClient email.Client, jwtAuth *jwtauth.JWTAuth) Service {
	return &service{
		logger:      logger,
		emailClient: emailClient,
		jwtAuth:     jwtAuth,
	}
}
