package email

import (
	"embed"

	"github.com/go-chi/jwtauth"
	"github.com/woningfinder/woningfinder/internal/customer"
	"github.com/woningfinder/woningfinder/pkg/email"
	"github.com/woningfinder/woningfinder/pkg/logging"
)

//go:embed templates/*
var emailTemplates embed.FS

// Service permits to handle the management of the email
type Service interface {
	SendWelcome(user *customer.User) error
	SendEmailConfirmationReminder(user *customer.User) error
	SendThankYou(user *customer.User) error
	SendFreeTrialReminder(user *customer.User) error
	SendLogin(user *customer.User) error
	SendWeeklyUpdate(user *customer.User, housingMatch []customer.HousingPreferencesMatch) error
	SendCorporationCredentialsMissing(user *customer.User) error
	SendCorporationCredentialsError(user *customer.User, corporationName string) error
	SendBye(user *customer.User) error

	ContactFormSubmission(content string) error
}

type service struct {
	logger      *logging.Logger
	emailClient email.Client
	jwtAuth     *jwtauth.JWTAuth
}

// NewService instantiate the email service
func NewService(logger *logging.Logger, emailClient email.Client, jwtAuth *jwtauth.JWTAuth) Service {
	return &service{
		logger:      logger,
		emailClient: emailClient,
		jwtAuth:     jwtAuth,
	}
}
