package cache

import (
	"time"

	"github.com/patrickmn/go-cache"
)

type CacheItf interface {
	Set(key string, value interface{}) error
	Get(key string) (interface{}, error)
}

func NewCache(cacheType string) CacheItf {
	switch cacheType {
	case "go-cache":
		return &GoCache{
			cache: cache.New(5*time.Minute, 10*time.Minute),
		}
	default:
		return &GoCache{
			cache: cache.New(5*time.Minute, 10*time.Minute),
		}
	}
}
