# Installation Guide

This guide provides instructions for installing and setting up the Go Caching Proxy.

## Prerequisites

- Go 1.21 or higher
- Git (for cloning the repository)

## Installing from Source

### 1. Clone the Repository

```bash
git clone https://github.com/user/go-caching-proxy.git
cd go-caching-proxy
```

### 2. Build the Application

```bash
go build -o caching-proxy ./cmd/caching-proxy
```

This will create an executable called `caching-proxy` in your current directory.

### 3. Verify Installation

Verify that the installation was successful by checking the version:

```bash
./caching-proxy --help
```

You should see the help information for the caching-proxy command.

## Using Docker

Alternatively, you can use Docker to run Go Caching Proxy without installing Go.

### 1. Building the Docker Image

```bash
docker build -t caching-proxy .
```

### 2. Using Docker Compose

The repository includes a `docker-compose.yml` file for easy deployment:

```bash
docker-compose up -d
```

This will start the caching proxy server on the port specified in the docker-compose.yml file.

## Next Steps

After installation, refer to the [Usage Guide](usage.md) for instructions on how to use the proxy server.
