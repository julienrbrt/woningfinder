package job

import (
	"github.com/woningfinder/woningfinder/internal/database"
	matcherService "github.com/woningfinder/woningfinder/internal/services/matcher"
	notificationsService "github.com/woningfinder/woningfinder/internal/services/notifications"
	userService "github.com/woningfinder/woningfinder/internal/services/user"
	"github.com/woningfinder/woningfinder/pkg/logging"
)

type Jobs struct {
	logger               *logging.Logger
	dbClient             database.DBClient
	redisClient          database.RedisClient
	userService          userService.Service
	matcherService       matcherService.Service
	notificationsService notificationsService.Service
}

func NewJobs(logger *logging.Logger, dbClient database.DBClient, redisClient database.RedisClient, userService userService.Service, matcherService matcherService.Service, notificationsService notificationsService.Service) *Jobs {
	return &Jobs{
		logger:               logger,
		dbClient:             dbClient,
		redisClient:          redisClient,
		userService:          userService,
		matcherService:       matcherService,
		notificationsService: notificationsService,
	}
}
