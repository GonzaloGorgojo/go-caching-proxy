package cache

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

type Cache struct {
	Key       string
	Body      []byte
	CreatedAt time.Time
	TTL       int64
	DB        *sql.DB
}

func (c *Cache) SetCache(key string, body []byte, ttl int) error {
	insertSQL := `INSERT INTO cache (key, body, created_at, ttl) VALUES (?, ?, ?, ?)`

	createdAt := time.Now().Format(time.DateTime)

	_, err := c.DB.Exec(insertSQL, key, body, createdAt, ttl)
	if err != nil {
		return err
	}

	log.Printf("Cache entry set for key: %s", key)
	return nil
}

func (c *Cache) GetCache(key string) (*Cache, error) {
	querySQL := `SELECT key, body, created_at, ttl FROM cache WHERE key = ?`

	var cacheEntry Cache
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

func (c *Cache) ClearCache() error {
	var count int
	getRowsQuery := `SELECT COUNT(key) FROM cache`

	err := c.DB.QueryRow(getRowsQuery).Scan(&count)
	if err != nil {
		log.Printf("Error counting cache: %s", err)
		return err
	}

	fmt.Printf("Cache status before clearing: %d entries\n", count)

	deleteQuery := `DELETE from cache`
	_, err = c.DB.Exec(deleteQuery)
	if err != nil {
		log.Printf("Error clearing cache: %s", err)
		return err
	}

	fmt.Println("Cache cleared entirely")
	return nil
}

func (c *Cache) ClearExpiredEntries() error {
	deleteExpiredQuery := `DELETE FROM cache WHERE DATETIME(created_at, '+' || ttl || ' minutes') <= DATETIME('now');`

	result, err := c.DB.Exec(deleteExpiredQuery)
	if err != nil {
		log.Printf("Error clearing expired cache entries: %s", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error retrieving number of deleted rows: %s", err)
		return err
	}

	log.Printf("Cleared %d expired cache entries.", rowsAffected)
	return nil

}

func (c *Cache) CleanUpService(interval time.Duration) {
	go func() {
		for {
			time.Sleep(interval)
			log.Printf("Running cleanup service every %v minutes\n", interval.Minutes())
			err := c.ClearExpiredEntries()
			if err != nil {
				log.Printf("Error during cleanup: %v", err)
			}
		}
	}()
}
