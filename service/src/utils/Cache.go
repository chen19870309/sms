package utils

import (
	"time"

	"github.com/patrickmn/go-cache"
)

var localCache *cache.Cache

func init() {
	localCache = cache.New(5*time.Minute, 30*time.Second)
}

func SetCache(k string, d interface{}, t time.Duration) {
	localCache.Set(k, d, t)
}

func GetCache(k string) (interface{}, bool) {
	return localCache.Get(k)
}

func DelCache(k string) {
	localCache.Delete(k)
}
