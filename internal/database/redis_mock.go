package database

import "github.com/go-redis/redis"

type RedisClientMock interface {
	Queue
	KeyStore
}

type redisClientMock struct {
	queueOutput    <-chan *redis.Message
	keyStoreOutput string
	err            error
}

func NewRedisClientMock(queueOutput <-chan *redis.Message, keyStoreOutput string, err error) RedisClientMock {
	return &redisClientMock{
		queueOutput:    queueOutput,
		keyStoreOutput: keyStoreOutput,
		err:            err,
	}
}

// Queue
func (c *redisClientMock) Publish(_ string, _ []byte) error {
	return c.err

}

func (c *redisClientMock) Subscribe(channelName string) (<-chan *redis.Message, error) {
	if c.err != nil {
		return nil, c.err
	}

	return c.queueOutput, nil
}

// KeyStore
func (c *redisClientMock) Get(key string) (string, error) {
	if c.err != nil {
		return "", c.err
	}

	return c.keyStoreOutput, nil
}

func (c *redisClientMock) Set(key string, value interface{}) error {
	return c.err
}
