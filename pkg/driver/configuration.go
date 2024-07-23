package driver

import (
	"fmt"
	"sync"
	"time"

	"github.com/adamnasrudin03/go-skeleton-chi/app/configs"
	"github.com/go-redis/redis"
)

var (
	cfg             = configs.GetInstance()
	logger          = Logger(cfg)
	redisClientOnce sync.Once
	redisClient     RedisClient
)

func Redis(config *configs.Configs) RedisClient {
	redisClientOnce.Do(func() {

		addrs := fmt.Sprintf("%s:%d", config.Redis.Host, config.Redis.Port)
		redisConn := redis.NewClient(&redis.Options{
			Addr:         addrs,
			Password:     config.Redis.Password,
			DB:           config.Redis.Database,
			PoolSize:     config.Redis.PoolSize,
			PoolTimeout:  time.Duration(config.Redis.PoolTimeout) * time.Second,
			MinIdleConns: config.Redis.MinIdleConn,
		})

		redisClient = NewRedis(redisConn)
	})

	return redisClient
}
