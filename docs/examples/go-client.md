# Go Client Example

This example demonstrates how to use the Go Caching Proxy from a Go application.

## Setup

First, ensure that the caching proxy is running:

```bash
./caching-proxy --port 3000 --origin https://jsonplaceholder.typicode.com
```

## Example Go Client

Create a file named `client.go` with the following content:

```go
package main

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

func main() {
	// Configure HTTP client
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Define proxy URL
	proxyURL := "http://localhost:3000"
	
	// First request (cache miss)
	fmt.Println("Making first request (cache miss)...")
	response1, headers1 := makeRequest(client, proxyURL+"/posts/1")
	fmt.Println("Response:", response1)
	fmt.Println("Cache status:", headers1.Get("X-Cache"))
	fmt.Println()

	// Second request (cache hit)
	fmt.Println("Making second request (cache hit)...")
	response2, headers2 := makeRequest(client, proxyURL+"/posts/1")
	fmt.Println("Response:", response2)
	fmt.Println("Cache status:", headers2.Get("X-Cache"))
	fmt.Println()

	// Different request (cache miss)
	fmt.Println("Making different request (cache miss)...")
	response3, headers3 := makeRequest(client, proxyURL+"/posts/2")
	fmt.Println("Response:", response3)
	fmt.Println("Cache status:", headers3.Get("X-Cache"))
}

func makeRequest(client *http.Client, url string) (string, http.Header) {
	// Create a new request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return "", nil
	}

	// Set headers
	req.Header.Set("Accept", "application/json")

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return "", nil
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return "", nil
	}

	return string(body), resp.Header
}
```

## Running the Example

Compile and run the example:

```bash
go run client.go
```

You should see output similar to:

```
Making first request (cache miss)...
Response: {
  "userId": 1,
  "id": 1,
  "title": "sunt aut facere repellat provident occaecati excepturi optio reprehenderit",
  "body": "quia et suscipit\nsuscipit recusandae consequuntur expedita et cum\nreprehenderit molestiae ut ut quas totam\nnostrum rerum est autem sunt rem eveniet architecto"
}
Cache status: MISS

Making second request (cache hit)...
Response: {
  "userId": 1,
  "id": 1,
  "title": "sunt aut facere repellat provident occaecati excepturi optio reprehenderit",
  "body": "quia et suscipit\nsuscipit recusandae consequuntur expedita et cum\nreprehenderit molestiae ut ut quas totam\nnostrum rerum est autem sunt rem eveniet architecto"
}
Cache status: HIT

Making different request (cache miss)...
Response: {
  "userId": 1,
  "id": 2,
  "title": "qui est esse",
  "body": "est rerum tempore vitae\nsequi sint nihil reprehenderit dolor beatae ea dolores neque\nfugiat blanditiis voluptate porro vel nihil molestiae ut reiciendis\nqui aperiam non debitis possimus qui neque nisi nulla"
}
Cache status: MISS
```

## Example Explanation

This example demonstrates:

1. Making HTTP requests through the caching proxy
2. Checking the cache status through the `X-Cache` header
3. Demonstrating that identical requests are served from cache (HIT)
4. Demonstrating that different requests result in cache misses (MISS)

## Making POST Requests

You can also make POST requests, which will bypass the cache:

```go
func makePostRequest(client *http.Client, url string, data string) (string, http.Header) {
	// Create a new request with a body
	req, err := http.NewRequest("POST", url, strings.NewReader(data))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return "", nil
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return "", nil
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return "", nil
	}

	return string(body), resp.Header
}
```

When using this function, you'll see the `X-Cache: BYPASS` header in the response.
