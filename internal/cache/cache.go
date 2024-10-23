package cache

import (
	"database/sql"
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

type Port struct {
	Port int
	DB   *sql.DB
}

func (p *Port) SetPort(port int) error {

	insertSQL := `INSERT or REPLACE INTO port (port) VALUES (?)`

	_, err := p.DB.Exec(insertSQL, port)
	if err != nil {
		return err
	}
	log.Printf("Port entry set for port: %d", port)
	return nil

}

func (p *Port) GetPort() (int, error) {
	var port int
	querySQL := `SELECT port FROM port ORDER BY id DESC LIMIT 1`

	err := p.DB.QueryRow(querySQL).Scan(&port)
	if err != nil {
		return 0, err
	}

	return port, nil
}

type DBCache struct {
	Key       string
	Body      []byte
	CreatedAt time.Time
	TTL       int64
	DB        *sql.DB
}

func (c *DBCache) SetCache(key string, body []byte, ttl int) error {
	insertSQL := `INSERT INTO cache (key, body, created_at, ttl) VALUES (?, ?, ?, ?)`

	createdAt := time.Now()

	_, err := c.DB.Exec(insertSQL, key, body, createdAt, ttl)
	if err != nil {
		return err
	}

	log.Printf("Cache entry set for key: %s", key)
	return nil
}

func (c *DBCache) GetCache(key string) (*DBCache, error) {
	querySQL := `SELECT key, body, created_at, ttl FROM cache WHERE key = ?`

	var cacheEntry DBCache
	err := c.DB.QueryRow(querySQL, key).Scan(&cacheEntry.Key, &cacheEntry.Body, &cacheEntry.CreatedAt, &cacheEntry.TTL)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("No cache entry found for key: %s", key)
			return nil, nil
		}
		return nil, err
	}

	return &cacheEntry, nil
}
