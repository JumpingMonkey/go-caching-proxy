# Benchmarking Guide

This guide provides instructions on how to benchmark the Go Caching Proxy to measure its performance characteristics.

## Performance Metrics

When benchmarking a caching proxy, several key metrics are important to measure:

1. **Throughput**: Requests per second that the proxy can handle
2. **Latency**: Response time for requests (both cache hits and misses)
3. **Cache Hit Ratio**: Percentage of requests served from cache
4. **Memory Usage**: RAM consumption under different loads

## Benchmarking Tools

### Using Apache Bench (ab)

[Apache Bench](https://httpd.apache.org/docs/2.4/programs/ab.html) is a simple tool for benchmarking HTTP servers.

#### Example: Testing Throughput and Latency

```bash
# Start the proxy server
./caching-proxy --port 3000 --origin http://dummyjson.com

# Run the benchmark (1000 requests, 10 concurrent)
ab -n 1000 -c 10 http://localhost:3000/products/1
```

This will output statistics including requests per second and response times.

#### Example: Comparing Cache Hit vs. Miss Performance

```bash
# Clear the cache first
./caching-proxy --clear-cache

# First run (cache miss)
ab -n 100 -c 10 http://localhost:3000/products/1

# Second run (cache hit)
ab -n 100 -c 10 http://localhost:3000/products/1
```

Compare the results to see the performance improvement with cached responses.

### Using hey

[hey](https://github.com/rakyll/hey) is another HTTP load generator that can be used for benchmarking.

```bash
# Install hey
go install github.com/rakyll/hey@latest

# Run benchmark
hey -n 1000 -c 50 http://localhost:3000/products/1
```

### Using Go's Built-in Benchmarking

You can also create Go benchmark tests for the proxy components.

Create a file named `proxy_bench_test.go` with content like:

```go
package proxy_test

import (
	"net/http"
	"testing"
)

func BenchmarkProxyGetCached(b *testing.B) {
	client := &http.Client{}
	
	// Warm up the cache
	req, _ := http.NewRequest("GET", "http://localhost:3000/products/1", nil)
	client.Do(req)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req, _ := http.NewRequest("GET", "http://localhost:3000/products/1", nil)
		client.Do(req)
	}
}
```

Run the benchmark with:

```bash
go test -bench=. -benchtime=10s
```

## Load Testing with Realistic Scenarios

For more realistic benchmarking, use tools like [Locust](https://locust.io/) or [Gatling](https://gatling.io/) to simulate real-world traffic patterns.

### Example Locust Test

Create a file named `locustfile.py`:

```python
from locust import HttpUser, task, between

class ProxyUser(HttpUser):
    wait_time = between(1, 3)
    
    @task(10)
    def get_popular_product(self):
        self.client.get("/products/1")
    
    @task(5)
    def get_random_products(self):
        for product_id in range(2, 6):
            self.client.get(f"/products/{product_id}")
    
    @task(1)
    def post_product(self):
        self.client.post("/products", json={
            "title": "Test Product",
            "description": "Test Description"
        })
```

Run Locust with:

```bash
locust -H http://localhost:3000
```

## Monitoring During Benchmarks

While running benchmarks, monitor system resources:

```bash
# Monitor CPU and memory usage
top -pid $(pgrep caching-proxy)

# For more detailed memory analysis
go tool pprof -alloc_space http://localhost:6060/debug/pprof/heap
```

## Analyzing Cache Performance

To analyze cache performance, add instrumentation to track cache hits and misses. You could modify the proxy code to log cache statistics or expose them via an HTTP endpoint.

## Benchmark Reporting

When reporting benchmark results, include:

1. Hardware specifications
2. Operating system
3. Go version
4. Test parameters (concurrent users, request rate)
5. Results (throughput, latency percentiles, cache hit ratio)

## Tips for Accurate Benchmarking

1. Run benchmarks multiple times and calculate averages
2. Ensure the origin server is not the bottleneck
3. Use realistic data sizes and request patterns
4. Test with different cache sizes and expiration settings
5. Run benchmarks from a separate machine to avoid resource contention
