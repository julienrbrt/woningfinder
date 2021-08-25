package job

import (
	"errors"

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

// markAsSent saves that an user has received an email reminder (using uuid representing the user and the email)
func (j *Jobs) markAsSent(uuid string) {
	if err := j.redisClient.Set(uuid, true); err != nil {
		j.logger.Sugar().Errorf("error when saving reminder state to redis: %w", err)
	}
}

// isAlreadySent check if a user already received an email reminder
func (j *Jobs) isAlreadySent(uuid string) bool {
	_, err := j.redisClient.Get(uuid)
	if err != nil {
		if !errors.Is(err, database.ErrRedisKeyNotFound) {
			j.logger.Sugar().Errorf("error when getting state reminder from redis: %w", err)
		}

		// does not have received reminder
		return false
	}

	return true
}
