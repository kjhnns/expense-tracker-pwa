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

	schema := `CREATE TABLE IF NOT EXISTS users (
        phone_number TEXT PRIMARY KEY,
        verified BOOLEAN NOT NULL DEFAULT 0,
        display_name TEXT,
        email TEXT,
        notify_by_sms BOOLEAN NOT NULL DEFAULT 1,
        notify_by_email BOOLEAN NOT NULL DEFAULT 1,
        payment_methods TEXT
    );

    CREATE TABLE IF NOT EXISTS groups (
        id TEXT PRIMARY KEY,
        name TEXT NOT NULL,
        created_by TEXT NOT NULL REFERENCES users(phone_number),
        default_currency TEXT NOT NULL,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP
    );

    CREATE TABLE IF NOT EXISTS group_members (
        group_id TEXT NOT NULL REFERENCES groups(id),
        phone_number TEXT NOT NULL REFERENCES users(phone_number),
        display_name_override TEXT,
        PRIMARY KEY (group_id, phone_number)
    );

    CREATE TABLE IF NOT EXISTS login_tokens (
        token TEXT PRIMARY KEY,
        phone_number TEXT NOT NULL REFERENCES users(phone_number),
        expires_at DATETIME NOT NULL
    );

    CREATE INDEX IF NOT EXISTS idx_group_members_group_id ON group_members(group_id);
    CREATE INDEX IF NOT EXISTS idx_group_members_phone_number ON group_members(phone_number);`

	if _, err := db.Exec(schema); err != nil {
		log.Fatalf("create tables: %v", err)
	}

	return db
}
