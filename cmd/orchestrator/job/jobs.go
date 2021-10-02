package job

import (
	"github.com/woningfinder/woningfinder/internal/database"
	emailService "github.com/woningfinder/woningfinder/internal/services/email"
	matcherService "github.com/woningfinder/woningfinder/internal/services/matcher"
	userService "github.com/woningfinder/woningfinder/internal/services/user"
	"github.com/woningfinder/woningfinder/pkg/logging"
)

type Jobs struct {
	logger         *logging.Logger
	dbClient       database.DBClient
	redisClient    database.RedisClient
	userService    userService.Service
	matcherService matcherService.Service
	emailService   emailService.Service
}

func NewJobs(logger *logging.Logger, dbClient database.DBClient, redisClient database.RedisClient, userService userService.Service, matcherService matcherService.Service, emailService emailService.Service) *Jobs {
	return &Jobs{
		logger:         logger,
		dbClient:       dbClient,
		redisClient:    redisClient,
		userService:    userService,
		matcherService: matcherService,
		emailService:   emailService,
	}
}
