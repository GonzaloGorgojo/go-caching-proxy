package cache

import (
	"fmt"
	"log"
	"sync"
	"time"
)

type CacheEntry struct {
	Body      []byte
	CreatedAt time.Time
	TTL       time.Duration
}

type Cache struct {
	Entries map[string]CacheEntry
	Mutex   sync.Mutex
}

func NewCache() *Cache {
	return &Cache{
		Entries: make(map[string]CacheEntry),
	}
}

func (c *Cache) ClearExpiredEntries() {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	clearedEntries := 0

	for key, entry := range c.Entries {
		if time.Since(entry.CreatedAt) >= entry.TTL {
			delete(c.Entries, key)
			clearedEntries++
			log.Printf("Cleared expired cache entry: %s\n", key)
		}
	}

	if clearedEntries == 0 {
		log.Println("No expired entries to clear.")
	}
}

func (c *Cache) CleanUpService(interval time.Duration) {
	go func() {
		for {
			time.Sleep(interval)
			log.Printf("Running cleanup service every %v minutes\n", interval.Minutes())
			c.ClearExpiredEntries()
		}
	}()
}

func (c *Cache) ClearCache() {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	fmt.Printf("Cache status before clearing: %d entries\n", len(c.Entries))

	if len(c.Entries) == 0 {
		fmt.Println("Cache is already empty")
	} else {
		fmt.Println("Current cache contents:")
		for key, entry := range c.Entries {
			fmt.Printf("Key: %s, Body length: %d, Created: %s, TTL: %s\n",
				key, len(entry.Body), entry.CreatedAt, entry.TTL)
		}
	}

	c.Entries = make(map[string]CacheEntry)

	fmt.Printf("\nCache status after clearing: %d entries\n", len(c.Entries))
	fmt.Println("Cache cleared entirely")
}
