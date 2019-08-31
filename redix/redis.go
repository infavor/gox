package redix

import (
	"github.com/go-redis/redis"
	"github.com/hetianyi/gox/logger"
	"time"
)

// InitRedisClient initializes redis client.
func InitRedisClient(address, pass string, poolSize int) *redis.Client {
	for {
		logger.Info("connecting to redis server")
		redisClient := redis.NewClient(&redis.Options{
			Addr:     address,
			Password: pass, // no password set
			DB:       0,    // use default DB
			PoolSize: poolSize,
		})
		if _, err := redisClient.Ping().Result(); err != nil {
			logger.Error(err)
			time.Sleep(time.Second * 10)
			logger.Info("try reconnecting to redis server")
			continue
		}
		logger.Info("successfully connected to redis server")
		return redisClient
	}
}
