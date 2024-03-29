package email

import (
	"embed"

	jwtauth "github.com/go-chi/jwtauth/v5"
	"github.com/julienrbrt/woningfinder/internal/corporation"
	"github.com/julienrbrt/woningfinder/internal/customer"
	"github.com/julienrbrt/woningfinder/pkg/downloader"
	"github.com/julienrbrt/woningfinder/pkg/email"
	"github.com/julienrbrt/woningfinder/pkg/logging"
)

//go:embed templates/*
var emailTemplates embed.FS

// Service permits to handle the management of the email
type Service interface {
	SendActivationEmail(user *customer.User) error
	SendEmailConfirmationReminder(user *customer.User) error
	SendThankYou(user *customer.User) error
	SendLogin(user *customer.User) error
	SendWeeklyUpdate(user *customer.User, matches []*customer.HousingPreferencesMatch) error
	SendCorporationCredentialsFirstTimeAdded(user *customer.User) error
	SendCorporationCredentialsMissing(user *customer.User) error
	SendCorporationCredentialsError(user *customer.User, corporationName string) error
	SendReactionFailure(user *customer.User, corporationName string, offers []corporation.Offer) error
	SendBye(user *customer.User) error

	SendWaitingListConfirmation(email, city string) error
	ContactFormSubmission(name, email, message string) error
}

type service struct {
	logger      *logging.Logger
	emailClient email.Client
	jwtAuth     *jwtauth.JWTAuth
	imgClient   downloader.Client
}

// NewService instantiate the email service
func NewService(logger *logging.Logger, emailClient email.Client, jwtAuth *jwtauth.JWTAuth, imgClient downloader.Client) Service {
	return &service{
		logger:      logger,
		emailClient: emailClient,
		jwtAuth:     jwtAuth,
		imgClient:   imgClient,
	}
}
