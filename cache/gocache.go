package cache

import (
	"errors"

	"github.com/patrickmn/go-cache"
	gocache "github.com/patrickmn/go-cache"
)

type GoCache struct {
	cache *gocache.Cache
}

func (gc *GoCache) Get(key string) (interface{}, error) {

	value, found := gc.cache.Get(key)
	if !found {
		return "", errors.New("not found")
	}

	return value, nil
}

func (gc *GoCache) Set(key string, value interface{}) error {
	gc.cache.Set(key, value, cache.DefaultExpiration)

	return nil
}
