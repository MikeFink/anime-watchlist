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
		`CREATE TABLE IF NOT EXISTS anime (
			id INTEGER PRIMARY KEY,
			anilist_id INTEGER UNIQUE,
			title TEXT NOT NULL,
			title_english TEXT,
			title_romaji TEXT,
			description TEXT,
			cover_image TEXT,
			banner_image TEXT,
			status TEXT,
			format TEXT,
			episodes INTEGER,
			duration INTEGER,
			season TEXT,
			season_year INTEGER,
			genres TEXT,
			score REAL,
			popularity INTEGER,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS watchlist (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			anime_id INTEGER,
			added_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (anime_id) REFERENCES anime (id)
		)`,
		`CREATE INDEX IF NOT EXISTS idx_anime_anilist_id ON anime(anilist_id)`,
		`CREATE INDEX IF NOT EXISTS idx_watchlist_anime_id ON watchlist(anime_id)`,
		`CREATE INDEX IF NOT EXISTS idx_anime_popularity ON anime(popularity DESC)`,
		`CREATE INDEX IF NOT EXISTS idx_anime_score ON anime(score DESC)`,
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