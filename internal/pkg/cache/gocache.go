package cache

import (
	"time"

	"github.com/patrickmn/go-cache"
)

type ICache interface {
	Get(key string) (val interface{}, isset bool)
	Set(key string, val interface{}, ttl time.Duration) error
	Delete(key string) error
	Flush() error
}

const (
	defaultExpiration = -1              // 永不失效
	cleanupInterval   = 5 * time.Minute // 清理间隔
)

type goCache struct {
	cache *cache.Cache
}

func NewGoCache() ICache {
	return &goCache{cache: cache.New(defaultExpiration, cleanupInterval)}
}

func (g *goCache) Get(key string) (val interface{}, isset bool) {
	return g.cache.Get(key)
}

func (g *goCache) Set(key string, val interface{}, ttl time.Duration) error {
	g.cache.Set(key, val, ttl)
	return nil
}

func (g *goCache) Delete(key string) error {
	g.cache.Delete(key)
	return nil
}

func (g *goCache) Flush() error {
	g.cache.Flush()
	return nil
}
