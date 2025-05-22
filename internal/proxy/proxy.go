package proxy

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"

	"github.com/user/go-caching-proxy/internal/cache"
)

// Start starts the caching proxy server
func Start(port int, originURL string) {
	// Parse the origin URL
	origin, err := url.Parse(originURL)
	if err != nil {
		log.Fatalf("Failed to parse origin URL: %v", err)
	}

	// Initialize the reverse proxy
	proxy := httputil.NewSingleHostReverseProxy(origin)

	// Create a custom director function to modify the request
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		req.Host = origin.Host
	}

	// Create a custom ModifyResponse function to cache responses
	proxy.ModifyResponse = func(resp *http.Response) error {
		// Skip caching for non-GET requests
		if resp.Request.Method != http.MethodGet {
			resp.Header.Set("X-Cache", "BYPASS")
			return nil
		}

		// Read the response body
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		resp.Body.Close()

		// Create a new ReadCloser from the body
		resp.Body = io.NopCloser(bytes.NewReader(body))

		// Store the response in the cache
		cacheKey := getCacheKey(resp.Request)
		headers := make(map[string][]string)
		for k, v := range resp.Header {
			headers[k] = v
		}

		cachedResp := cache.Response{
			StatusCode: resp.StatusCode,
			Headers:    headers,
			Body:       body,
		}
		cache.Set(cacheKey, cachedResp)

		// Set the X-Cache header
		resp.Header.Set("X-Cache", "MISS")

		return nil
	}

	// Create a custom handler for the proxy
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)

		// Skip cache for non-GET requests
		if r.Method != http.MethodGet {
			proxy.ServeHTTP(w, r)
			return
		}

		// Try to get the response from the cache
		cacheKey := getCacheKey(r)
		cachedResp, found := cache.Get(cacheKey)
		if found {
			// Set the status code
			w.WriteHeader(cachedResp.StatusCode)

			// Set the headers
			for k, values := range cachedResp.Headers {
				for _, v := range values {
					w.Header().Add(k, v)
				}
			}

			// Set the X-Cache header
			w.Header().Set("X-Cache", "HIT")

			// Write the body
			_, err := w.Write(cachedResp.Body)
			if err != nil {
				log.Printf("Failed to write cached response: %v", err)
			}
			return
		}

		// If not found in cache, forward the request to the origin server
		proxy.ServeHTTP(w, r)
	})

	// Start the server
	server := &http.Server{
		Addr:    ":" + strconv.Itoa(port),
		Handler: handler,
	}

	log.Fatal(server.ListenAndServe())
}

// getCacheKey generates a unique key for the request
func getCacheKey(r *http.Request) string {
	// Create a key based on the request method, URL, and headers
	key := r.Method + "|" + r.URL.String()

	// Add important headers to the key
	for _, h := range []string{"Accept", "Accept-Encoding"} {
		if v := r.Header.Get(h); v != "" {
			key += "|" + h + ":" + v
		}
	}

	// Hash the key to create a fixed-length string
	hash := md5.Sum([]byte(key))
	return hex.EncodeToString(hash[:])
}
