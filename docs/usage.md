# Usage Guide

This guide explains how to use the Go Caching Proxy server after installation.

## Basic Usage

Start the caching proxy server by specifying the port and origin server URL.

### Command Line Options

- `--port <number>`: The port on which the proxy server will listen (required)
- `--origin <url>`: The URL of the origin server to which requests will be forwarded (required)
- `--clear-cache`: Clear the cache and exit

### Basic Example

```bash
./caching-proxy --port 3000 --origin http://dummyjson.com
```

This starts a caching proxy server on port 3000 that forwards requests to http://dummyjson.com.

## Docker Usage

If you're using Docker, you can run the proxy as follows:

```bash
docker run -p 3000:3000 caching-proxy --port 3000 --origin http://dummyjson.com
```

## Clearing the Cache

To clear the cache:

```bash
./caching-proxy --clear-cache
```

## How It Works

1. When you make a GET request to the proxy server (e.g., `http://localhost:3000/products`), it checks if the response is already cached.
2. If the response is not in the cache, it forwards the request to the origin server, caches the response, and returns it with an `X-Cache: MISS` header.
3. If the response is in the cache, it returns the cached response with an `X-Cache: HIT` header.
4. Non-GET requests (POST, PUT, DELETE) are always forwarded to the origin server and not cached.

## Cache Headers

The proxy adds the following custom headers to responses:

- `X-Cache: HIT` - When the response was served from cache
- `X-Cache: MISS` - When the response was fetched from the origin server
- `X-Cache: BYPASS` - When caching was bypassed (for non-GET requests)

## Example Workflow

1. Start the proxy server:
   ```bash
   ./caching-proxy --port 3000 --origin http://dummyjson.com
   ```

2. Make a request to the proxy:
   ```bash
   curl -i http://localhost:3000/products/1
   ```
   The response will include the `X-Cache: MISS` header.

3. Make the same request again:
   ```bash
   curl -i http://localhost:3000/products/1
   ```
   This time, the response will include the `X-Cache: HIT` header, indicating it was served from cache.
