# Architecture

This document provides an overview of the Go Caching Proxy architecture, explaining its components and how they interact.

## System Overview

Go Caching Proxy is built as a lightweight HTTP proxy server that caches responses to improve performance. The system consists of the following main components:

1. Command-line interface (CLI)
2. Proxy server
3. Caching system

## Component Diagram

```
┌─────────────────┐      ┌─────────────────┐      ┌─────────────────┐
│                 │      │                 │      │                 │
│   HTTP Client   │──────▶     Proxy      │──────▶  Origin Server  │
│                 │      │                 │      │                 │
└─────────────────┘      └────────┬────────┘      └─────────────────┘
                                  │
                                  │
                          ┌───────▼────────┐
                          │                │
                          │     Cache      │
                          │                │
                          └────────────────┘
```

## Components

### Command-line Interface (CLI)

The CLI is implemented using the Cobra library and is responsible for parsing command-line arguments and starting the proxy server. It accepts parameters such as the port number, origin server URL, and a flag to clear the cache.

**Key files:**
- `/cmd/caching-proxy/main.go`

### Proxy Server

The proxy server is built on top of Go's standard `net/http` and `net/http/httputil` packages. It forwards HTTP requests to the origin server and handles caching logic. The proxy component is responsible for:

1. Forwarding requests to the origin server
2. Intercepting responses
3. Checking if responses are cacheable
4. Generating cache keys
5. Serving cached responses when available

**Key files:**
- `/internal/proxy/proxy.go`

### Caching System

The caching system stores HTTP responses and retrieves them when needed. It uses the `github.com/patrickmn/go-cache` package for in-memory caching with expiration. The cache component provides:

1. Storage for HTTP responses
2. Key-based retrieval
3. Automatic expiration
4. Cache clearing functionality

**Key files:**
- `/internal/cache/cache.go`

## Request Flow

1. Client sends a request to the proxy server
2. Proxy server generates a cache key for the request
3. Proxy checks if the response is already in cache
   - If found, returns the cached response with `X-Cache: HIT` header
   - If not found, forwards the request to the origin server
4. When receiving response from origin server:
   - For cacheable responses (GET requests), stores in cache and returns with `X-Cache: MISS` header
   - For non-cacheable responses, returns with `X-Cache: BYPASS` header

## Cache Key Generation

Cache keys are generated based on:
- Request method
- URL
- Important headers like `Accept` and `Accept-Encoding`

This ensures that different variations of the same request (e.g., requesting JSON vs. XML) are cached separately.

## Caching Policies

- Only GET requests are cached
- Default cache expiration time is 5 minutes
- Cache is automatically purged every 10 minutes

## Future Improvements

- Add support for cache control headers
- Implement disk-based caching for persistence
- Add HTTP-based cache management API
- Support for more customization options (cache size, expiration time, etc.)
- Add metrics and monitoring capabilities
