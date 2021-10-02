package database

import "errors"

// SetUUID saves an UUID in redis
func (r *redisClient) SetUUID(uuid string) {
	if err := r.Set(uuid, true); err != nil {
		r.logger.Sugar().Errorf("error when saving reminder state to redis: %w", err)
	}
}

// HasUUID check if an UUID is aleady stored in redis
func (r *redisClient) HasUUID(uuid string) bool {
	_, err := r.Get(uuid)
	if err != nil {
		if !errors.Is(err, ErrRedisKeyNotFound) {
			r.logger.Sugar().Errorf("error when getting state reminder from redis: %w", err)
		}

		// does not have UUID
		return false
	}

	return true
}
