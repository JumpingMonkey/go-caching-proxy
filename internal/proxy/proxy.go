package proxy

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"

	"github.com/user/go-caching-proxy/internal/cache"
)

// Custom ResponseWriter to intercept and modify responses
type proxyResponseWriter struct {
	http.ResponseWriter
	isNonGetRequest bool
	isCacheMiss     bool
}

// WriteHeader intercepts the WriteHeader call to add our custom headers
func (w *proxyResponseWriter) WriteHeader(status int) {
	if w.isNonGetRequest {
		// Remove any existing X-Cache header and set our BYPASS header
		w.Header().Del("X-Cache")
		w.Header().Set("X-Cache", "BYPASS")
	} else if w.isCacheMiss {
		// Remove any existing X-Cache header and set our MISS header
		w.Header().Del("X-Cache")
		w.Header().Set("X-Cache", "MISS")
	}
	w.ResponseWriter.WriteHeader(status)
}

// handleCachedResponse handles a request with a cached response
func handleCachedResponse(w http.ResponseWriter, r *http.Request, cacheKey string) bool {
	cachedResp, found := cache.Get(cacheKey)
	if !found {
		return false
	}

	// Set the X-Cache header first so we have control
	w.Header().Set("X-Cache", "HIT")
	
	// Set the headers, but skip X-Cache which we want to control
	for k, values := range cachedResp.Headers {
		// Skip headers we want to control
		if k == "X-Cache" {
			continue
		}
		for _, v := range values {
			w.Header().Add(k, v)
		}
	}

	// Set the X-Cache header again to make sure it's not overridden
	w.Header().Set("X-Cache", "HIT")

	// Set the status code
	w.WriteHeader(cachedResp.StatusCode)

	// Write the body
	_, err := w.Write(cachedResp.Body)
	if err != nil {
		log.Printf("Failed to write cached response: %v", err)
		return false
	}
	
	return true
}

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
			// Remove any existing X-Cache header from origin to avoid confusion
			resp.Header.Del("X-Cache")
			// Set our custom header
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

		// Remove any existing X-Cache header from origin
		resp.Header.Del("X-Cache")
		// Set the X-Cache header
		resp.Header.Set("X-Cache", "MISS")

		return nil
	}

	// Create a custom handler for the proxy
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)

		// Skip cache for non-GET requests
		if r.Method != http.MethodGet {
			// Create a custom ResponseWriter to intercept the response
			proxyWriter := &proxyResponseWriter{
				ResponseWriter: w,
				isNonGetRequest: true,
			}
			proxy.ServeHTTP(proxyWriter, r)
			return
		}

		// Try to get the response from the cache
		cacheKey := getCacheKey(r)
		
		// Try to handle the request with a cached response
		if handled := handleCachedResponse(w, r, cacheKey); handled {
			log.Printf("Cache HIT for %s", r.URL.Path)
			return
		} else {
			log.Printf("Cache MISS for %s", r.URL.Path)
			// Create a custom ResponseWriter to intercept the response
			proxyWriter := &proxyResponseWriter{
				ResponseWriter: w,
				isCacheMiss: true,
			}
			proxy.ServeHTTP(proxyWriter, r)
		}
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
	// Create a key based on the request method, path, and query only (not full URL)
	// This ensures consistency between client request and origin server request
	pathAndQuery := r.URL.Path
	if r.URL.RawQuery != "" {
		pathAndQuery += "?" + r.URL.RawQuery
	}
	
	key := r.Method + "|" + pathAndQuery

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
