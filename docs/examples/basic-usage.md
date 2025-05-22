# Basic Usage Example

This example demonstrates how to use the Go Caching Proxy with curl to cache API responses.

## Starting the Proxy

First, start the caching proxy server, pointing it to a public API:

```bash
./caching-proxy --port 3000 --origin https://jsonplaceholder.typicode.com
```

You should see output similar to:

```
Starting caching proxy server on port 3000, forwarding to https://jsonplaceholder.typicode.com
```

## Making Requests

### First Request (Cache Miss)

Make a request to the proxy server:

```bash
curl -i http://localhost:3000/posts/1
```

You should see a response similar to:

```
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
X-Cache: MISS
...

{
  "userId": 1,
  "id": 1,
  "title": "sunt aut facere repellat provident occaecati excepturi optio reprehenderit",
  "body": "quia et suscipit\nsuscipit recusandae consequuntur expedita et cum\nreprehenderit molestiae ut ut quas totam\nnostrum rerum est autem sunt rem eveniet architecto"
}
```

Note the `X-Cache: MISS` header, indicating that the response was not from cache.

### Second Request (Cache Hit)

Make the same request again:

```bash
curl -i http://localhost:3000/posts/1
```

This time, you should see:

```
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
X-Cache: HIT
...

{
  "userId": 1,
  "id": 1,
  "title": "sunt aut facere repellat provident occaecati excepturi optio reprehenderit",
  "body": "quia et suscipit\nsuscipit recusandae consequuntur expedita et cum\nreprehenderit molestiae ut ut quas totam\nnostrum rerum est autem sunt rem eveniet architecto"
}
```

Note the `X-Cache: HIT` header, indicating that the response was served from cache.

## Non-GET Requests

Make a POST request:

```bash
curl -i -X POST -H "Content-Type: application/json" -d '{"title":"foo","body":"bar","userId":1}' http://localhost:3000/posts
```

You should see a response with:

```
X-Cache: BYPASS
```

This indicates that the request was forwarded to the origin server without caching.

## Clearing the Cache

To clear the cache, run:

```bash
./caching-proxy --clear-cache
```

After clearing the cache, the next request will be a cache miss, even if it was previously cached.
