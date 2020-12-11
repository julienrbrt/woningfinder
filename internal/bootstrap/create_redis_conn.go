package bootstrap

import (
	"fmt"
	"log"
	"os"

	"github.com/go-redis/redis"
)

//RDB stores the redis connection
var RDB *redis.Client

// InitRedis create a connection to WoningFinder Redis storage
func InitRedis() error {
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	redisPassword := os.Getenv("REDIS_PASSWORD")

	opt, err := redis.ParseURL(fmt.Sprintf("rediss://default:%s@%s:%s/0", redisPassword, redisHost, redisPort))
	if err != nil {
		return fmt.Errorf("error connecting to redis with host: %s", redisHost)
	}

	rdb := redis.NewClient(opt)

	_, err = rdb.Ping().Result()
	if err != nil {
		return fmt.Errorf("error connecting to redis with host: %s", redisHost)
	}

	RDB = rdb
	if RDB != nil {
		log.Println("successfully connected to redis")
	}

	return nil
}
