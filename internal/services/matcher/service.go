package matcher

import (
	"context"
	"fmt"

	"github.com/julienrbrt/woningfinder/internal/corporation"
	"github.com/julienrbrt/woningfinder/internal/corporation/connector"
	"github.com/julienrbrt/woningfinder/internal/customer/matcher"
	"github.com/julienrbrt/woningfinder/internal/database"
	corporationService "github.com/julienrbrt/woningfinder/internal/services/corporation"
	emailService "github.com/julienrbrt/woningfinder/internal/services/email"
	userService "github.com/julienrbrt/woningfinder/internal/services/user"
	"github.com/julienrbrt/woningfinder/pkg/downloader"
	"github.com/julienrbrt/woningfinder/pkg/logging"
)

// Service permits to handle the persistence of matcher
type Service interface {
	Match(ctx context.Context, offers corporation.Offers) error
}

type service struct {
	logger             *logging.Logger
	dbClient           database.DBClient
	userService        userService.Service
	emailService       emailService.Service
	corporationService corporationService.Service
	imgClient          downloader.Client
	matcher            matcher.Matcher
	connectorProvider  connector.ConnectorProvider
}

// NewService instantiate the matcher service
func NewService(logger *logging.Logger, dbClient database.DBClient, userService userService.Service, emailService emailService.Service, corporationService corporationService.Service, imgClient downloader.Client, matcher matcher.Matcher, connectorProvider connector.ConnectorProvider) Service {
	return &service{
		logger:             logger,
		dbClient:           dbClient,
		userService:        userService,
		emailService:       emailService,
		corporationService: corporationService,
		imgClient:          imgClient,
		matcher:            matcher,
		connectorProvider:  connectorProvider,
	}
}

// MatcherCounter holds the user reactions to offers
type MatcherCounter struct {
	UUID         string `pg:",pk"`
	Reacted      bool
	FailureCount uint
}

func (s *service) setMatch(m *MatcherCounter) error {
	if _, err := s.dbClient.Conn().
		Model(m).
		OnConflict("(UUID) DO UPDATE").
		Insert(); err != nil {
		return fmt.Errorf("error when updating match counter for %s: %w", m.UUID, err)
	}

	return nil
}

func (s *service) getMatch(uuid string) (*MatcherCounter, error) {
	match := &MatcherCounter{UUID: uuid}
	if err := s.dbClient.Conn().
		Model(match).
		Where("UUID = ?", uuid).
		Select(); err != nil {
		return nil, fmt.Errorf("error when getting match counter for %s: %w", uuid, err)
	}

	return match, nil
}
