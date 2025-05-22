# Go Caching Proxy Documentation

Welcome to the documentation for Go Caching Proxy, a lightweight HTTP proxy server that caches responses for improved performance.

## Table of Contents

- [Overview](#overview)
- [Installation](installation.md)
- [Usage](usage.md)
- [API Reference](api-reference.md)
- [Architecture](architecture.md)
- [Examples](examples/README.md)
- [Benchmarking](benchmarking.md)
- [Troubleshooting](troubleshooting.md)
- [Contributing](contributing.md)

## Overview

Go Caching Proxy is a command-line tool that starts a caching proxy server. It forwards HTTP requests to an origin server and caches the responses. Subsequent identical requests will return the cached response instead of forwarding to the server, improving performance and reducing load on the origin server.

### Key Features

- Forward HTTP requests to a specified origin server
- Cache GET responses for improved performance
- Custom headers indicating cache status (HIT/MISS)
- Command-line interface for easy configuration
- Cache clearing functionality
- Docker support for containerized deployment

## Quick Start

### Installation

```bash
# Clone the repository
git clone https://github.com/user/go-caching-proxy.git
cd go-caching-proxy

# Build the application
go build -o caching-proxy ./cmd/caching-proxy
```

### Basic Usage

```bash
# Start the proxy server
./caching-proxy --port 3000 --origin http://dummyjson.com

# Make a request to the proxy
curl -i http://localhost:3000/products/1
```

## Documentation Structure

- [Installation Guide](installation.md): Detailed instructions for installing the proxy
- [Usage Guide](usage.md): How to use the proxy with various configuration options
- [API Reference](api-reference.md): Detailed description of the proxy's HTTP interface and internal APIs
- [Architecture](architecture.md): Overview of the system architecture and components
- [Examples](examples/README.md): Practical examples for different use cases
- [Benchmarking](benchmarking.md): How to measure the proxy's performance
- [Troubleshooting](troubleshooting.md): Solutions for common issues
- [Contributing](contributing.md): Guidelines for contributing to the project
