package cache

import (
	"sync"
	"time"

	"github.com/patrickmn/go-cache"
)

var (
	// Create a cache with a default expiration time of 5 minutes, and which
	// purges expired items every 10 minutes
	memoryCache *cache.Cache
	once        sync.Once
)

// Response represents a cached HTTP response
type Response struct {
	StatusCode int
	Headers    map[string][]string
	Body       []byte
}

// Initialize initializes the cache
func Initialize() {
	once.Do(func() {
		memoryCache = cache.New(5*time.Minute, 10*time.Minute)
	})
}

// Set stores a response in the cache
func Set(key string, response Response) {
	Initialize()
	memoryCache.Set(key, response, cache.DefaultExpiration)
}

// Get retrieves a response from the cache
func Get(key string) (Response, bool) {
	Initialize()
	value, found := memoryCache.Get(key)
	if !found {
		return Response{}, false
	}

	return value.(Response), true
}

// Clear removes all items from the cache
func Clear() {
	Initialize()
	memoryCache.Flush()
}
