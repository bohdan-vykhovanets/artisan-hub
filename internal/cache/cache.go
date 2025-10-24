package cache

import (
	"time"

	"github.com/patrickmn/go-cache"
)

var AppCache *cache.Cache
var DefaultExpiration time.Duration = cache.DefaultExpiration

func InitCache() {
	AppCache = cache.New(5*time.Minute, 10*time.Minute)
}
