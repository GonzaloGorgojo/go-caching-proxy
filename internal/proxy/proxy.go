package proxy

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/gonzalogorgojo/go-caching-proxy/internal/cache"
)

func ProxyHandler(target string, c *cache.Cache) (http.HandlerFunc, error) {
	targetURL, err := url.Parse(target)
	if err != nil {
		return nil, fmt.Errorf("failed to parse target URL: %w", err)
	}

	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	return func(w http.ResponseWriter, r *http.Request) {
		cacheKey := r.URL.String()

		c.Mutex.Lock()
		defer c.Mutex.Unlock()

		cachedEntry, found := c.Entries[cacheKey]

		if found {
			fmt.Println("Using Cache response")

			w.Header().Add("X-Cache", "Hit")
			w.WriteHeader(http.StatusOK)
			w.Write(cachedEntry.Body)
			return
		}

		proxy.ModifyResponse = func(r *http.Response) error {

			body, err := io.ReadAll(r.Body)
			if err != nil {
				return fmt.Errorf("failed to read response body: %w", err)
			}

			r.Body.Close()
			r.Body = io.NopCloser(bytes.NewReader(body))

			cacheEntry := cache.CacheEntry{
				Body:      body,
				CreatedAt: time.Now(),
				TTL:       60 * time.Second,
			}

			c.Entries[cacheKey] = cacheEntry

			r.Header.Add("X-Cache", "Miss")
			return nil
		}

		fmt.Printf("Forwarding request to: %s\n", targetURL)
		proxy.ServeHTTP(w, r)
	}, nil
}
