package database

import (
	"fmt"
	"time"

	"anime-watchlist/backend/domain"
)

type WatchlistRepository struct {
	db *Database
}

func NewWatchlistRepository(db *Database) *WatchlistRepository {
	return &WatchlistRepository{db: db}
}

func (r *WatchlistRepository) GetWatchlist() ([]domain.WatchlistItem, error) {
	query := `
		SELECT id, anilist_id, added_at
		FROM watchlist
		ORDER BY added_at DESC
	`

	rows, err := r.db.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query watchlist: %w", err)
	}
	defer rows.Close()

	var items []domain.WatchlistItem
	for rows.Next() {
		var item domain.WatchlistItem
		err := rows.Scan(&item.ID, &item.AnilistID, &item.AddedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan watchlist item: %w", err)
		}
		items = append(items, item)
	}

	return items, nil
}

func (r *WatchlistRepository) AddToWatchlist(anilistID int) error {
	query := `
		INSERT INTO watchlist (anilist_id, added_at)
		VALUES (?, ?)
	`

	_, err := r.db.DB.Exec(query, anilistID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to add to watchlist: %w", err)
	}

	return nil
}

func (r *WatchlistRepository) RemoveFromWatchlist(anilistID int) error {
	query := `
		DELETE FROM watchlist
		WHERE anilist_id = ?
	`

	result, err := r.db.DB.Exec(query, anilistID)
	if err != nil {
		return fmt.Errorf("failed to remove from watchlist: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("anime not in watchlist")
	}

	return nil
}

func (r *WatchlistRepository) IsInWatchlist(anilistID int) (bool, error) {
	query := `
		SELECT COUNT(*) FROM watchlist
		WHERE anilist_id = ?
	`

	var count int
	err := r.db.DB.QueryRow(query, anilistID).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check watchlist: %w", err)
	}

	return count > 0, nil
}

func (r *WatchlistRepository) GetWatchlistCount() (int, error) {
	query := `SELECT COUNT(*) FROM watchlist`

	var count int
	err := r.db.DB.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get watchlist count: %w", err)
	}

	return count, nil
} 