package database

import (
	"fmt"

	"github.com/go-redis/redis"
)

// Queue defines how the different Redis pub-sub is used (as queue)
type Queue interface {
	Publish(channelName string, data []byte) error
	Subscribe(channelName string) (<-chan *redis.Message, error)
}

func (r *redisClient) Publish(channelName string, data []byte) error {
	result := r.rdb.Publish(channelName, data)
	if result.Err() != nil {
		return fmt.Errorf("error publishing %v to channel %s: %w", string(data), channelName, result.Err())
	}

	if result.Val() == 0 {
		r.logger.Sugar().Warnf("⚠️ warning successfuly published %v to channel %s: no one received it", string(data), channelName)
	}

	return nil
}

func (r *redisClient) Subscribe(channelName string) (<-chan *redis.Message, error) {
	pubsub := r.rdb.Subscribe(channelName)
	// defer pubsub.Close()

	// wait for confirmation that subscription is created before doing anything.
	if _, err := pubsub.Receive(); err != nil {
		return nil, fmt.Errorf("error subscribing to channel: %w", err)
	}

	//channel which receives messages.
	return pubsub.Channel(), nil
}
