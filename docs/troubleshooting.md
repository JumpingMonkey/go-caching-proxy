# Troubleshooting Guide

This guide helps you diagnose and resolve common issues with the Go Caching Proxy.

## Connection Issues

### Proxy Server Won't Start

**Symptoms:**
- Error message when starting the proxy
- Process exits immediately

**Possible Causes and Solutions:**

1. **Port already in use**
   ```
Error: listen tcp :3000: bind: address already in use
   ```
   
   **Solution:** Choose a different port or stop the process using the current port:
   ```bash
   # Find the process using the port
   lsof -i :3000
   
   # Kill the process
   kill <PID>
   ```

2. **Invalid origin URL**
   ```
Failed to parse origin URL: parse "http//example.com": invalid URL
   ```
   
   **Solution:** Ensure the origin URL is correctly formatted with the proper scheme (http:// or https://).

### Cannot Connect to Proxy

**Symptoms:**
- Connection refused when trying to access the proxy

**Possible Causes and Solutions:**

1. **Proxy not running**
   
   **Solution:** Ensure the proxy is running. Check the process:
   ```bash
   ps aux | grep caching-proxy
   ```

2. **Firewall blocking the connection**
   
   **Solution:** Check firewall settings and ensure the port is open.

3. **Trying to connect from another machine**
   
   **Solution:** By default, the proxy binds to 0.0.0.0, but ensure your network allows the connection.

## Caching Issues

### Responses Not Being Cached

**Symptoms:**
- Every request shows `X-Cache: MISS`

**Possible Causes and Solutions:**

1. **Non-GET requests**
   
   **Solution:** Only GET requests are cached. Other methods (POST, PUT, DELETE) always bypass the cache.

2. **Unique request parameters**
   
   **Solution:** Check if the requests are actually identical. Query parameters, different Accept headers, or other varying headers can create different cache keys.

3. **Cache expiration**
   
   **Solution:** The default cache expiration is 5 minutes. After this time, items are evicted from the cache.

### Cache Not Clearing

**Symptoms:**
- `--clear-cache` doesn't seem to work

**Possible Causes and Solutions:**

1. **Using different instances**
   
   **Solution:** Make sure you're clearing the cache on the same instance that's running the proxy.

2. **Permission issues**
   
   **Solution:** Ensure you have the necessary permissions to modify the cache.

## Performance Issues

### Slow Response Times

**Symptoms:**
- Responses take longer than expected, even with cache hits

**Possible Causes and Solutions:**

1. **Origin server is slow**
   
   **Solution:** For cache misses, the proxy depends on the origin server's performance. Check if the origin server is responding slowly.

2. **High load**
   
   **Solution:** Check system resources (CPU, memory) to ensure the proxy has enough resources.

3. **Network latency**
   
   **Solution:** Check network conditions between client, proxy, and origin server.

### High Memory Usage

**Symptoms:**
- The proxy process is using more memory than expected

**Possible Causes and Solutions:**

1. **Large responses being cached**
   
   **Solution:** The proxy caches entire responses in memory. Large responses will consume more memory.

2. **Many unique requests**
   
   **Solution:** Each unique request (based on method, URL, and some headers) gets a separate cache entry.

## Docker Issues

### Container Exits Immediately

**Symptoms:**
- Docker container starts and then exits

**Possible Causes and Solutions:**

1. **Command-line arguments missing**
   
   **Solution:** Ensure you're providing the required `--port` and `--origin` flags.

2. **Port conflict in container**
   
   **Solution:** Check if the port is already in use within the container.

### Cannot Access Proxy from Host

**Symptoms:**
- Proxy is running in Docker but not accessible from the host

**Possible Causes and Solutions:**

1. **Port mapping missing**
   
   **Solution:** Ensure you've mapped the container port to the host with `-p` flag:
   ```bash
   docker run -p 3000:3000 caching-proxy --port 3000 --origin http://example.com
   ```

2. **Container networking issue**
   
   **Solution:** Check Docker network settings.

## Debugging

### Enabling Debug Output

To get more information about what's happening, you can add debugging output by modifying the source code in `internal/proxy/proxy.go`:

```go
// Add this near the top
import (
    // existing imports
    "os"
)

// Add this to the Start function
verbose := os.Getenv("PROXY_DEBUG") == "1"

// Add debug logs throughout the code
if verbose {
    log.Printf("Debug: Cache key generated: %s", cacheKey)
}
```

Then run the proxy with debugging enabled:

```bash
PROXY_DEBUG=1 ./caching-proxy --port 3000 --origin http://example.com
```

### Inspecting Cache Contents

Currently, there's no built-in way to inspect the cache contents. However, you could add an endpoint to expose this information by modifying the source code.

## Getting Help

If you're still experiencing issues:

1. Check the GitHub repository for known issues
2. Open a new issue on GitHub with detailed information about your problem
3. Include logs, error messages, and steps to reproduce the issue
