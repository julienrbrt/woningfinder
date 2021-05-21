package database

import (
	"fmt"

	"github.com/go-redis/redis"
	"github.com/woningfinder/woningfinder/pkg/logging"
)

type RedisClient interface {
	Queue
	KeyStore
}

type redisClient struct {
	logger *logging.Logger
	rdb    *redis.Client
}

// NewRedisClient creates a connection to WoningFinder Redis storage
func NewRedisClient(logger *logging.Logger, host, port, password string) (RedisClient, error) {
	options, err := redis.ParseURL(fmt.Sprintf("rediss://default:%s@%s:%s/0", password, host, port))
	if err != nil {
		return nil, fmt.Errorf("error connecting to redis with host: %s", host)
	}

	// retry query 3 times (instead of not)
	options.MaxRetries = 3

	rdb := redis.NewClient(options)

	_, err = rdb.Ping().Result()
	if err != nil {
		return nil, fmt.Errorf("error connecting to redis with host: %s", host)
	}

	if rdb != nil {
		logger.Sugar().Info("successfully connected to redis 🎉")
	}

	return &redisClient{
		logger: logger,
		rdb:    rdb,
	}, nil
}
