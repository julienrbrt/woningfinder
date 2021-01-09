package database

import (
	"fmt"

	"github.com/go-redis/redis"
)

// PubSubOffers defines on which channel the corporation offers are sent via redis
const PubSubOffers = "offers"

// PubSubPayment defines on which channel the payment confirmation are sent via redis
const PubSubPayment = "stripe"

// Queue defines how the different Redis pub-sub is used (as queue)
type Queue interface {
	Publish(channelName string, data []byte) error
	Subscribe(channelName string) (<-chan *redis.Message, error)
}

func (r *redisClient) Publish(channelName string, data []byte) error {
	if err := r.rdb.Publish(channelName, data).Err(); err != nil {
		return fmt.Errorf("error publishing %v to channel %s: %w", data, channelName, err)
	}

	return nil
}

func (r *redisClient) Subscribe(channelName string) (<-chan *redis.Message, error) {
	pubsub := r.rdb.Subscribe(channelName)
	defer pubsub.Close()

	// wait for confirmation that subscription is created before doing anything.
	if _, err := pubsub.Receive(); err != nil {
		return nil, fmt.Errorf("error subscribing to channel: %w", err)
	}

	//channel which receives messages.
	return pubsub.Channel(), nil
}
