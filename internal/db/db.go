package db

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func InitDB() *sql.DB {
	if db != nil {
		return db
	}

	var err error

	db, err = sql.Open("sqlite3", "cache_data.db")
	if err != nil {
		log.Fatal(err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	log.Println("Database initiated correctly")

	createUserTable(db)
	return db
}

func createUserTable(db *sql.DB) {
	createSQLTable := `
    CREATE TABLE IF NOT EXISTS cache (
        "id" INTEGER PRIMARY KEY AUTOINCREMENT,
        "key" TEXT NOT NULL UNIQUE,
        "body" BLOB NOT NULL,
		"created_at" DATETIME NOT NULL,
		"ttl" INTEGER NOT NULL
    );`

	_, err := db.Exec(createSQLTable)

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Cache table created succesfully.")
}
