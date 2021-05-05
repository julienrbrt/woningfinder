package job

import (
	"github.com/woningfinder/woningfinder/internal/database"
	matcherService "github.com/woningfinder/woningfinder/internal/services/matcher"
	notificationService "github.com/woningfinder/woningfinder/internal/services/notification"
	userService "github.com/woningfinder/woningfinder/internal/services/user"
	"github.com/woningfinder/woningfinder/pkg/logging"
)

type Jobs struct {
	logger              *logging.Logger
	dbClient            database.DBClient
	redisClient         database.RedisClient
	userService         userService.Service
	matcherService      matcherService.Service
	notificationService notificationService.Service
}

func NewJobs(logger *logging.Logger, dbClient database.DBClient, redisClient database.RedisClient, userService userService.Service, matcherService matcherService.Service, notificationService notificationService.Service) *Jobs {
	return &Jobs{
		logger:              logger,
		dbClient:            dbClient,
		redisClient:         redisClient,
		userService:         userService,
		matcherService:      matcherService,
		notificationService: notificationService,
	}
}
