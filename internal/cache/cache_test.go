package cache

import (
	"testing"
	"time"

	"github.com/patrickmn/go-cache"
)

func TestCacheInitialize(t *testing.T) {
	// Initialize the cache
	Initialize()

	// Check that the memory cache is not nil
	if memoryCache == nil {
		t.Error("Expected memoryCache to be initialized, got nil")
	}
}

func TestCacheSetGet(t *testing.T) {
	// Initialize the cache
	Initialize()

	// Create a test response
	testResp := Response{
		StatusCode: 200,
		Headers: map[string][]string{
			"Content-Type": {"application/json"},
		},
		Body: []byte(`{"test": "data"}`),
	}

	// Store the response in the cache
	Set("test-key", testResp)

	// Retrieve the response from the cache
	cachedResp, found := Get("test-key")
	if !found {
		t.Error("Expected to find response in cache, but not found")
	}

	// Check that the cached response matches the original
	if cachedResp.StatusCode != testResp.StatusCode {
		t.Errorf("Expected StatusCode %d, got %d", testResp.StatusCode, cachedResp.StatusCode)
	}

	if cachedResp.Headers["Content-Type"][0] != testResp.Headers["Content-Type"][0] {
		t.Errorf("Expected Content-Type %s, got %s", testResp.Headers["Content-Type"][0], cachedResp.Headers["Content-Type"][0])
	}

	if string(cachedResp.Body) != string(testResp.Body) {
		t.Errorf("Expected Body %s, got %s", string(testResp.Body), string(cachedResp.Body))
	}
}

func TestCacheExpiration(t *testing.T) {
	// Create a custom cache with a short expiration for testing
	// directly set the memory cache instead of using Initialize()
	memoryCache = cache.New(100*time.Millisecond, 200*time.Millisecond)

	// Create a test response
	testResp := Response{
		StatusCode: 200,
		Headers:    map[string][]string{},
		Body:       []byte(`{"test": "expiration"}`),
	}

	// Store the response in the cache
	Set("expire-key", testResp)

	// Verify the response is in the cache
	_, found := Get("expire-key")
	if !found {
		t.Error("Expected to find response in cache immediately after setting, but not found")
	}

	// Wait for the cache to expire
	time.Sleep(150 * time.Millisecond)

	// Verify the response is no longer in the cache
	_, found = Get("expire-key")
	if found {
		t.Error("Expected response to be expired from cache, but it was found")
	}
}

func TestCacheClear(t *testing.T) {
	// Initialize the cache
	Initialize()

	// Create a test response
	testResp := Response{
		StatusCode: 200,
		Headers:    map[string][]string{},
		Body:       []byte(`{"test": "clear"}`),
	}

	// Store the response in the cache
	Set("clear-key", testResp)

	// Verify the response is in the cache
	_, found := Get("clear-key")
	if !found {
		t.Error("Expected to find response in cache, but not found")
	}

	// Clear the cache
	Clear()

	// Verify the response is no longer in the cache
	_, found = Get("clear-key")
	if found {
		t.Error("Expected cache to be empty after Clear(), but found response")
	}
}
