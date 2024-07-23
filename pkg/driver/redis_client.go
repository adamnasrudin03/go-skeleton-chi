package driver

import (
	"time"

	help "github.com/adamnasrudin03/go-helpers"
	"github.com/go-redis/redis"
)

type RedisClient interface {
	Del(key string) error
	Set(key string, value interface{}, expDur time.Duration) error
	Get(key string) (string, error)
}

type redisCtx struct {
	redisClient redis.Cmdable
}

func NewRedis(redisClient redis.Cmdable) RedisClient {
	return &redisCtx{
		redisClient: redisClient,
	}
}

func (c *redisCtx) Del(key string) error {
	err := c.redisClient.Del(key).Err()
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (c *redisCtx) Get(key string) (string, error) {
	data, err := c.redisClient.Get(key).Result()
	if err != nil {
		logger.Error(err)
		return "", err
	}
	return data, nil
}

func (c *redisCtx) Set(key string, value interface{}, expDur time.Duration) error {
	payload, err := help.SafeJsonMarshal(value)
	if err != nil {
		logger.Error(err)
		return err
	}

	err = c.redisClient.Set(key, payload, expDur).Err()
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (c *redisCtx) Keys(key string) ([]string, error) {
	data, err := c.redisClient.Keys(key).Result()
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	return data, nil
}
