# Docker Compose Example

This example demonstrates how to use Docker Compose to set up the Go Caching Proxy with another service.

## Docker Compose Configuration

Create a file named `docker-compose.example.yml` with the following content:

```yaml
version: '3'

services:
  # The caching proxy service
  proxy:
    build:
      context: .
      dockerfile: Dockerfile
    ports:
      - "3000:3000"
    command: --port 3000 --origin http://api:3001
    networks:
      - app-network
    restart: unless-stopped
    depends_on:
      - api

  # Mock API service (using json-server as an example)
  api:
    image: vimagick/json-server
    ports:
      - "3001:3001"
    volumes:
      - ./test/data:/data
    command: -H 0.0.0.0 -p 3001 -w /data/db.json
    networks:
      - app-network
    restart: unless-stopped

networks:
  app-network:
    driver: bridge
```

## Test Data

Create a test data file at `test/data/db.json` with the following content:

```json
{
  "posts": [
    { "id": 1, "title": "Post 1", "author": "Author 1", "content": "Content 1" },
    { "id": 2, "title": "Post 2", "author": "Author 2", "content": "Content 2" },
    { "id": 3, "title": "Post 3", "author": "Author 3", "content": "Content 3" }
  ],
  "comments": [
    { "id": 1, "postId": 1, "text": "Comment 1" },
    { "id": 2, "postId": 1, "text": "Comment 2" },
    { "id": 3, "postId": 2, "text": "Comment 3" }
  ],
  "profile": {
    "name": "Example Profile"
  }
}
```

## Starting the Services

Start the services using Docker Compose:

```bash
docker-compose -f docker-compose.example.yml up -d
```

This will start two containers:
1. A caching proxy on port 3000
2. A mock JSON API on port 3001

## Testing the Setup

### Direct API Access

Access the API directly:

```bash
curl -i http://localhost:3001/posts/1
```

### Through the Caching Proxy

Access the API through the proxy:

```bash
curl -i http://localhost:3000/posts/1
```

The first request will show `X-Cache: MISS` and subsequent identical requests will show `X-Cache: HIT`.

## Scaling the Setup

This setup can be extended for more complex scenarios:

### Adding a Load Balancer

```yaml
services:
  # Load balancer
  nginx:
    image: nginx:latest
    ports:
      - "80:80"
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf:ro
    depends_on:
      - proxy
    networks:
      - app-network
```

With a corresponding `nginx.conf` file:

```nginx
events {
    worker_connections 1024;
}

http {
    upstream backend {
        server proxy:3000;
    }

    server {
        listen 80;
        
        location / {
            proxy_pass http://backend;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
        }
    }
}
```

## Monitoring with Docker Stats

Monitor the containers:

```bash
docker stats
```

## Stopping the Services

Stop and remove the containers:

```bash
docker-compose -f docker-compose.example.yml down
```

## Benefits of This Setup

1. **Isolation**: Each service runs in its own container
2. **Networking**: Services can communicate over a private network
3. **Caching**: The proxy caches API responses for improved performance
4. **Scalability**: Easy to scale by adding more services
5. **Reproducibility**: Environment can be easily recreated on any machine with Docker
