package repository

import (
	"fmt"
	"goapi/pkg/module/auth/core/ports"
	"goapi/pkg/util/cache"
	"time"
)

const (
	prefix = "TOKEN::" // กำหนด prefix ของ key สำหรับ token
)

type tokenRepositiry struct {
	cache *cache.RedisClient
}

func NewTokenRepository(cache *cache.RedisClient) ports.TokenRepository {
	return &tokenRepositiry{
		cache: cache,
	}
}

func (r tokenRepositiry) SetToken(tokenId string, data map[string]any, duration time.Duration) error {
	return r.cache.Set(fmt.Sprintf("%s%s", prefix, tokenId), data, duration)
}

func (r tokenRepositiry) GetToken(tokenId string) (string, error) {
	return r.cache.Get(fmt.Sprintf("%s%s", prefix, tokenId))
}

func (r tokenRepositiry) DeleteToken(tokenId string) (int64, error) {
	return r.cache.Delete(fmt.Sprintf("%s%s", prefix, tokenId))
}
