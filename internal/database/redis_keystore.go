package database

import (
	"context"
	"errors"
	"fmt"

	redis "github.com/go-redis/redis/v8"
)

// KeyStore defines how the redis is used (as a keystore)
type KeyStore interface {
	Get(key string) (string, error)
	Set(key string, value interface{}) error
}

// ErrRedisKeyNotFound is an error returned when a key does not have a value
var ErrRedisKeyNotFound = errors.New("key not found")

func (r *redisClient) Get(key string) (string, error) {
	value, err := r.rdb.Get(context.Background(), key).Result()
	if err != nil {
		if err != redis.Nil {
			return "", fmt.Errorf("error when getting key %s from redis: %w", key, err)
		}

		return "", fmt.Errorf("error finding key %s: %w", key, ErrRedisKeyNotFound)
	}

	return value, nil
}

func (r *redisClient) Set(key string, value interface{}) error {
	return r.rdb.Set(context.Background(), key, value, 0).Err()
}
