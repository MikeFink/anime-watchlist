package database

import (
	"database/sql"
	"fmt"
	"log"
	"path/filepath"

	_ "modernc.org/sqlite"
)

type Database struct {
	DB *sql.DB
}

func New(dbPath string) (*Database, error) {
	dir := filepath.Dir(dbPath)
	if err := ensureDir(dir); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %w", err)
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	database := &Database{DB: db}
	if err := database.migrate(); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return database, nil
}

func (d *Database) Close() error {
	return d.DB.Close()
}

func (d *Database) migrate() error {
	queries := []string{
		`DROP TABLE IF EXISTS watchlist`,
		`DROP TABLE IF EXISTS anime`,
		`CREATE TABLE IF NOT EXISTS watchlist (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			anilist_id INTEGER UNIQUE NOT NULL,
			added_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE INDEX IF NOT EXISTS idx_watchlist_anilist_id ON watchlist(anilist_id)`,
		`CREATE INDEX IF NOT EXISTS idx_watchlist_added_at ON watchlist(added_at DESC)`,
	}

	for _, query := range queries {
		if _, err := d.DB.Exec(query); err != nil {
			return fmt.Errorf("failed to execute migration: %w", err)
		}
	}

	log.Println("Database migration completed successfully")
	return nil
}

func ensureDir(dir string) error {
	if dir == "." || dir == "/" {
		return nil
	}
	return nil
} 