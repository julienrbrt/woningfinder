package database

type RedisClientMock interface {
	KeyStore
	Queue
}

type redisClientMock struct {
	RedisClient
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
	// this error cannot happen in set to it has been defined in the mock for the Get
	// to bypass
	if c.err == ErrRedisKeyNotFound {
		return nil
	}

	return c.err
}

func (c *redisClientMock) HasUUID(uuid string) bool {
	return c.err == nil
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
