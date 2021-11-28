package database

import (
	"context"
	"fmt"

	redis "github.com/go-redis/redis/v8"
	"github.com/julienrbrt/woningfinder/pkg/logging"
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
func NewRedisClient(logger *logging.Logger, redisURL string) (RedisClient, error) {
	options, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("error connecting to redis with host: %w", err)
	}

	// retry query 3 times (instead of not)
	options.MaxRetries = 3

	rdb := redis.NewClient(options)

	_, err = rdb.Ping(context.Background()).Result()
	if err != nil {
		return nil, fmt.Errorf("error connecting to redis: %w", err)
	}

	if rdb != nil {
		logger.Info("successfully connected to redis 🎉")
	}

	return &redisClient{
		logger: logger,
		rdb:    rdb,
	}, nil
}
