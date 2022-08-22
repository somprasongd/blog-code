package database

import (
	"errors"
	"fmt"
	"goapi/pkg/config"
	"goapi/pkg/util/cache"

	"github.com/go-redis/redis/v9"
)

func NewRedisClient(conf *config.Config) (*cache.RedisClient, error) {
	if len(conf.Redis.Host) == 0 || conf.Redis.Port == 0 {
		return nil, nil
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%v", conf.Redis.Host, conf.Redis.Port),
		Password: conf.Redis.Password, // empty mean no password set
		DB:       conf.Redis.Database, // 0 is default DB
	})

	redisCache := cache.New(rdb)

	if !redisCache.Ping() {
		return nil, errors.New("redis connection failed")
	}

	return redisCache, nil
}
