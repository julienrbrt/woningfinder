package database

import (
	"errors"

	"go.uber.org/zap"
)

// SetUUID saves an UUID in redis
func (r *redisClient) SetUUID(uuid string) {
	if err := r.Set(uuid, true); err != nil {
		r.logger.Error("error when saving uuid state to redis", zap.Error(err))
	}
}

// HasUUID check if an UUID is aleady stored in redis
func (r *redisClient) HasUUID(uuid string) bool {
	_, err := r.Get(uuid)
	if err != nil {
		if !errors.Is(err, ErrRedisKeyNotFound) {
			r.logger.Error("error when getting uuid state from redis", zap.Error(err))
		}

		// does not have UUID
		return false
	}

	return true
}
