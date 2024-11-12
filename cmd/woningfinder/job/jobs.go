package job

import (
	"github.com/julienrbrt/woningfinder/internal/database"
	emailService "github.com/julienrbrt/woningfinder/internal/services/email"
	matcherService "github.com/julienrbrt/woningfinder/internal/services/matcher"
	userService "github.com/julienrbrt/woningfinder/internal/services/user"
	"github.com/julienrbrt/woningfinder/pkg/logging"
)

type Jobs struct {
	logger         *logging.Logger
	dbClient       database.DBClient
	userService    userService.Service
	matcherService matcherService.Service
	emailService   emailService.Service
}

func NewJobs(logger *logging.Logger, dbClient database.DBClient, userService userService.Service, matcherService matcherService.Service, emailService emailService.Service) *Jobs {
	return &Jobs{
		logger:         logger,
		dbClient:       dbClient,
		userService:    userService,
		matcherService: matcherService,
		emailService:   emailService,
	}
}
