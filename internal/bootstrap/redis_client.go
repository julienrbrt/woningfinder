package bootstrap

import (
	"github.com/woningfinder/woningfinder/internal/config"
	"github.com/woningfinder/woningfinder/internal/database"
	"go.uber.org/zap"
)

func CreateRedisClient(logger *zap.Logger) database.RedisClient {
	client, err := database.NewRedisClient(logger, config.MustGetString("REDIS_HOST"), config.MustGetString("REDIS_PORT"), config.MustGetString("REDIS_PASSWORD"))
	if err != nil {
		logger.Sugar().Fatal(err)
	}

	return client
}
