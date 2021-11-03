package bootstrap

import (
	"github.com/woningfinder/woningfinder/internal/database"
	"github.com/woningfinder/woningfinder/pkg/config"
	"github.com/woningfinder/woningfinder/pkg/logging"
	"go.uber.org/zap"
)

// CreateRedisClient creates a connection to redis and provides its client
func CreateRedisClient(logger *logging.Logger) database.RedisClient {
	client, err := database.NewRedisClient(logger, config.MustGetString("REDIS_HOST"), config.MustGetString("REDIS_PORT"), config.MustGetString("REDIS_PASSWORD"))
	if err != nil {
		logger.Fatal("error when creating redis client", zap.Error(err))
	}

	return client
}
