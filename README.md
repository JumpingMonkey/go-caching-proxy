# Go Caching Proxy

A command-line tool that starts a caching proxy server, which forwards requests to an origin server and caches the responses. Subsequent identical requests will return the cached response instead of forwarding to the server.

## Features

- Forward HTTP requests to a specified origin server
- Cache GET responses for improved performance
- Custom headers indicating cache status (HIT/MISS)
- Command-line interface for easy configuration
- Cache clearing functionality
- Docker support for containerized deployment

## Installation

### Using Go

If you have Go installed:

```bash
# Clone the repository
git clone https://github.com/user/go-caching-proxy.git
cd go-caching-proxy

# Build the application
go build -o caching-proxy ./cmd/caching-proxy
```

### Using Docker

```bash
# Build the Docker image
docker build -t caching-proxy .

# Or use docker-compose
docker-compose up -d
```

## Usage

### Basic Usage

Start the caching proxy server by specifying the port and origin server URL:

```bash
./caching-proxy --port <number> --origin <url>
```

Example:

```bash
./caching-proxy --port 3000 --origin http://dummyjson.com
```

This starts a caching proxy server on port 3000 that forwards requests to http://dummyjson.com.

### Docker Usage

```bash
docker run -p 3000:3000 caching-proxy --port 3000 --origin http://dummyjson.com
```

### Clear Cache

To clear the cache:

```bash
./caching-proxy --clear-cache
```

## How It Works

1. When you make a GET request to the proxy server (e.g., `http://localhost:3000/products`), it checks if the response is already cached.
2. If the response is not in the cache, it forwards the request to the origin server, caches the response, and returns it with an `X-Cache: MISS` header.
3. If the response is in the cache, it returns the cached response with an `X-Cache: HIT` header.
4. Non-GET requests (POST, PUT, DELETE) are always forwarded to the origin server and not cached.

## Cache Key Generation

The cache key is generated based on:
- Request method
- Path and query parameters (not the full URL with domain)
- Important headers like `Accept` and `Accept-Encoding`

This approach ensures:
- Different variations of the same request (e.g., requesting JSON vs. XML) are cached separately
- Consistent caching behavior regardless of which domain or hostname is used to access the proxy
- Better cache hit ratios by using domain-independent keys

## Documentation

Comprehensive documentation is available in the `docs/` directory:

- Installation and usage guides
- API reference
- Architecture overview
- Examples and troubleshooting

To view the documentation, browse the files in the `docs/` directory, starting with `docs/README.md`.

## Testing

The project includes comprehensive test suites for both the cache and proxy components. These tests verify the correct functionality of the caching mechanism, response headers, and proxy behavior.

### Running Tests with Docker

The easiest way to run tests is using Docker:

```bash
# Build the test Docker image
docker build -t caching-proxy-test -f Dockerfile.test .

# Run the tests
docker run --rm caching-proxy-test
```

### Test Coverage

The test suite covers the following functionality:

#### Cache Tests
- Cache initialization and configuration
- Setting and retrieving cached responses
- Cache expiration behavior
- Cache clearing functionality

#### Proxy Tests
- Cache key generation (ensuring domain-independence)
- Proxy response writer headers (HIT/MISS/BYPASS)
- Cached response handling
- Proper header management

## Development

### Project Structure

```
go-caching-proxy/
├── cmd/              # Command-line interface code
├── internal/         # Internal application code
│   ├── cache/        # Cache implementation
│   └── proxy/        # Proxy server implementation
├── docs/             # API documentation
├── Dockerfile        # Docker configuration
├── docker-compose.yml # Docker compose configuration
├── go.mod            # Go module definition
└── README.md         # Project documentation
```

## License

MIT
