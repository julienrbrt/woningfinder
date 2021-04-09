package notifications

import (
	"github.com/go-chi/jwtauth"
	"github.com/woningfinder/woningfinder/internal/entity"
	"github.com/woningfinder/woningfinder/pkg/email"
	"github.com/woningfinder/woningfinder/pkg/logging"
)

// Service permits to handle the management of the notifications
type Service interface {
	SendWelcome(user *entity.User) error
	SendWeeklyUpdate(user *entity.User, housingMatch []entity.HousingPreferencesMatch) error
	SendCorporationCredentialsError(user *entity.User, corporationName string) error
	SendBye(user *entity.User) error
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
