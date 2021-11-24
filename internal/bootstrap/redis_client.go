package bootstrap

import (
	"github.com/julienrbrt/woningfinder/internal/database"
	"github.com/julienrbrt/woningfinder/pkg/config"
	"github.com/julienrbrt/woningfinder/pkg/logging"
	"go.uber.org/zap"
)

// CreateRedisClient creates a connection to redis and provides its client
func CreateRedisClient(logger *logging.Logger) database.RedisClient {
	client, err := database.NewRedisClient(logger, config.MustGetString("REDIS_URL"))
	if err != nil {
		logger.Fatal("error when creating redis client", zap.Error(err))
	}

	return client
}
