package database

type RedisClientMock interface {
	KeyStore
	Queue
}

type redisClientMock struct {
	keyStoreOutput string
	queueOutput    []string
	err            error
}

func NewRedisClientMock(keyStoreOutput string, queueOutput []string, err error) RedisClientMock {
	return &redisClientMock{
		keyStoreOutput: keyStoreOutput,
		queueOutput:    queueOutput,
		err:            err,
	}
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

// Queue
func (c *redisClientMock) Push(listName string, data []byte) error {
	return c.err
}

func (c *redisClientMock) BPop(listName string) ([]string, error) {
	if c.err != nil {
		return nil, c.err
	}

	return c.queueOutput, nil
}
