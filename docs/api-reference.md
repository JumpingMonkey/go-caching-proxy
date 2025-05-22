# API Reference

## HTTP Proxy Interface

The Go Caching Proxy provides an HTTP proxy interface that forwards requests to the origin server and caches responses.

### Request Format

Requests should be made to the proxy server with the path of the resource you want to access on the origin server.

```
GET http://<proxy-host>:<proxy-port>/<path>
```

Example:
```
GET http://localhost:3000/products/1
```

This will forward the request to the origin server at `<origin-url>/products/1`.

### Response Format

Responses from the proxy server include the original response from the origin server, plus additional headers indicating cache status.

#### Cache Headers

- `X-Cache: HIT` - Response was served from cache
- `X-Cache: MISS` - Response was fetched from the origin server
- `X-Cache: BYPASS` - Caching was bypassed (for non-GET requests)

### Caching Behavior

- Only GET requests are cached
- Cache keys are generated based on the request method, URL, and important headers (Accept, Accept-Encoding)
- Default cache expiration time is 5 minutes
- Cache is automatically purged every 10 minutes

## Cache Management

### Clearing the Cache

The cache can be cleared using the `--clear-cache` command line option:

```bash
./caching-proxy --clear-cache
```

There is currently no HTTP API for clearing the cache.

## Internal Package API

### Cache Package

```go
// Initialize initializes the cache
func Initialize()

// Set stores a response in the cache
func Set(key string, response Response)

// Get retrieves a response from the cache
func Get(key string) (Response, bool)

// Clear removes all items from the cache
func Clear()
```

### Proxy Package

```go
// Start starts the caching proxy server
func Start(port int, originURL string)

// getCacheKey generates a unique key for the request
func getCacheKey(r *http.Request) string
```
