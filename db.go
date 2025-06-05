package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

// initDB opens the SQLite database and ensures required tables exist.
func initDB(path string) *sql.DB {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		log.Fatalf("open db: %v", err)
	}

	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS groups (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        phone_numbers TEXT NOT NULL,
        name TEXT NOT NULL,
        created_by TEXT,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        default_currency TEXT
    )`); err != nil {
		log.Fatalf("create table: %v", err)
	}

	return db
}
