package cache

import (
	"fmt"
	"time"

	"github.com/patrickmn/go-cache"
)

var Cache = cache.New(5*time.Minute, 5*time.Minute)

const (
	GithubMessageCacheKey = "github_message"
)

func SetCache(key string, data interface{}) bool {
	fmt.Println("SET CACHE! KEY: ", key, " DATA: ", data)
	Cache.Set(key, data, 1*time.Minute)
	return true
}
func GetCache(key string) (string, bool) {
	var (
		data  string
		found bool
	)
	result, found := Cache.Get(key)
	if found {
		data = result.(string)
		fmt.Println("CACHE FOUND! KEY", key, " DATA: ", data)
	}
	return data, found
}
