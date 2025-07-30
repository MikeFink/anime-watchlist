package database

import (
	"database/sql"
	"fmt"
	"strings"

	"anime-watchlist/backend/domain"
)

type AnimeRepository struct {
	db *Database
}

func NewAnimeRepository(db *Database) *AnimeRepository {
	return &AnimeRepository{db: db}
}

func (r *AnimeRepository) GetAll(filter domain.AnimeFilter) ([]domain.Anime, error) {
	if err := filter.Validate(); err != nil {
		return nil, err
	}

	query := `
		SELECT 
			a.id, a.anilist_id, a.title, a.title_english, a.title_romaji,
			a.description, a.cover_image, a.banner_image, a.status, a.format,
			a.episodes, a.duration, a.season, a.season_year, a.genres,
			a.score, a.popularity, a.created_at, a.updated_at,
			CASE WHEN w.anime_id IS NOT NULL THEN 1 ELSE 0 END as is_watching
		FROM anime a
		LEFT JOIN watchlist w ON a.id = w.anime_id
	`

	args := []interface{}{}
	conditions := []string{}

	if filter.Search != "" {
		conditions = append(conditions, "(a.title LIKE ? OR a.title_english LIKE ?)")
		searchTerm := "%" + filter.Search + "%"
		args = append(args, searchTerm, searchTerm)
	}

	if filter.Status != "" {
		conditions = append(conditions, "a.status = ?")
		args = append(args, filter.Status)
	}

	if filter.Season != "" {
		conditions = append(conditions, "a.season = ?")
		args = append(args, filter.Season)
	}

	if filter.Year > 0 {
		conditions = append(conditions, "a.season_year = ?")
		args = append(args, filter.Year)
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += " ORDER BY a.popularity DESC, a.score DESC"
	query += " LIMIT ? OFFSET ?"

	offset := (filter.Page - 1) * filter.PageSize
	args = append(args, filter.PageSize, offset)

	rows, err := r.db.DB.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query anime: %w", err)
	}
	defer rows.Close()

	var animes []domain.Anime
	for rows.Next() {
		var anime domain.Anime
		err := rows.Scan(
			&anime.ID, &anime.AnilistID, &anime.Title, &anime.TitleEnglish, &anime.TitleRomaji,
			&anime.Description, &anime.CoverImage, &anime.BannerImage, &anime.Status, &anime.Format,
			&anime.Episodes, &anime.Duration, &anime.Season, &anime.SeasonYear, &anime.Genres,
			&anime.Score, &anime.Popularity, &anime.CreatedAt, &anime.UpdatedAt, &anime.IsWatching,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan anime: %w", err)
		}
		animes = append(animes, anime)
	}

	return animes, nil
}

func (r *AnimeRepository) GetByID(id int) (*domain.Anime, error) {
	query := `
		SELECT 
			a.id, a.anilist_id, a.title, a.title_english, a.title_romaji,
			a.description, a.cover_image, a.banner_image, a.status, a.format,
			a.episodes, a.duration, a.season, a.season_year, a.genres,
			a.score, a.popularity, a.created_at, a.updated_at,
			CASE WHEN w.anime_id IS NOT NULL THEN 1 ELSE 0 END as is_watching
		FROM anime a
		LEFT JOIN watchlist w ON a.id = w.anime_id
		WHERE a.id = ?
	`

	var anime domain.Anime
	err := r.db.DB.QueryRow(query, id).Scan(
		&anime.ID, &anime.AnilistID, &anime.Title, &anime.TitleEnglish, &anime.TitleRomaji,
		&anime.Description, &anime.CoverImage, &anime.BannerImage, &anime.Status, &anime.Format,
		&anime.Episodes, &anime.Duration, &anime.Season, &anime.SeasonYear, &anime.Genres,
		&anime.Score, &anime.Popularity, &anime.CreatedAt, &anime.UpdatedAt, &anime.IsWatching,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("anime not found: %d", id)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get anime: %w", err)
	}

	return &anime, nil
}

func (r *AnimeRepository) GetWatchlist() ([]domain.Anime, error) {
	query := `
		SELECT 
			a.id, a.anilist_id, a.title, a.title_english, a.title_romaji,
			a.description, a.cover_image, a.banner_image, a.status, a.format,
			a.episodes, a.duration, a.season, a.season_year, a.genres,
			a.score, a.popularity, a.created_at, a.updated_at,
			1 as is_watching
		FROM anime a
		INNER JOIN watchlist w ON a.id = w.anime_id
		ORDER BY w.added_at DESC
	`

	rows, err := r.db.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query watchlist: %w", err)
	}
	defer rows.Close()

	var animes []domain.Anime
	for rows.Next() {
		var anime domain.Anime
		err := rows.Scan(
			&anime.ID, &anime.AnilistID, &anime.Title, &anime.TitleEnglish, &anime.TitleRomaji,
			&anime.Description, &anime.CoverImage, &anime.BannerImage, &anime.Status, &anime.Format,
			&anime.Episodes, &anime.Duration, &anime.Season, &anime.SeasonYear, &anime.Genres,
			&anime.Score, &anime.Popularity, &anime.CreatedAt, &anime.UpdatedAt, &anime.IsWatching,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan anime: %w", err)
		}
		animes = append(animes, anime)
	}

	return animes, nil
}

func (r *AnimeRepository) AddToWatchlist(animeID int) error {
	tx, err := r.db.DB.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	var existingID int
	err = tx.QueryRow("SELECT id FROM anime WHERE id = ?", animeID).Scan(&existingID)
	if err == sql.ErrNoRows {
		return fmt.Errorf("anime not found: %d", animeID)
	}
	if err != nil {
		return fmt.Errorf("failed to check anime existence: %w", err)
	}

	var watchlistID int
	err = tx.QueryRow("SELECT id FROM watchlist WHERE anime_id = ?", animeID).Scan(&watchlistID)
	if err == nil {
		return fmt.Errorf("anime already in watchlist: %d", animeID)
	}
	if err != sql.ErrNoRows {
		return fmt.Errorf("failed to check watchlist: %w", err)
	}

	_, err = tx.Exec("INSERT INTO watchlist (anime_id) VALUES (?)", animeID)
	if err != nil {
		return fmt.Errorf("failed to add to watchlist: %w", err)
	}

	return tx.Commit()
}

func (r *AnimeRepository) RemoveFromWatchlist(animeID int) error {
	result, err := r.db.DB.Exec("DELETE FROM watchlist WHERE anime_id = ?", animeID)
	if err != nil {
		return fmt.Errorf("failed to remove from watchlist: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("anime not in watchlist: %d", animeID)
	}

	return nil
}

func (r *AnimeRepository) Upsert(anime *domain.Anime) error {
	if err := anime.Validate(); err != nil {
		return err
	}

	query := `
		INSERT OR REPLACE INTO anime (
			anilist_id, title, title_english, title_romaji, description,
			cover_image, banner_image, status, format, episodes, duration,
			season, season_year, genres, score, popularity, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP)
	`

	_, err := r.db.DB.Exec(query,
		anime.AnilistID, anime.Title, anime.TitleEnglish, anime.TitleRomaji, anime.Description,
		anime.CoverImage, anime.BannerImage, anime.Status, anime.Format, anime.Episodes, anime.Duration,
		anime.Season, anime.SeasonYear, anime.Genres, anime.Score, anime.Popularity,
	)
	if err != nil {
		return fmt.Errorf("failed to upsert anime: %w", err)
	}

	return nil
}

func (r *AnimeRepository) BulkUpsert(animes []domain.Anime) error {
	if len(animes) == 0 {
		return nil
	}

	tx, err := r.db.DB.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`
		INSERT OR REPLACE INTO anime (
			anilist_id, title, title_english, title_romaji, description,
			cover_image, banner_image, status, format, episodes, duration,
			season, season_year, genres, score, popularity, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP)
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	for _, anime := range animes {
		if err := anime.Validate(); err != nil {
			return fmt.Errorf("invalid anime: %w", err)
		}

		_, err := stmt.Exec(
			anime.AnilistID, anime.Title, anime.TitleEnglish, anime.TitleRomaji, anime.Description,
			anime.CoverImage, anime.BannerImage, anime.Status, anime.Format, anime.Episodes, anime.Duration,
			anime.Season, anime.SeasonYear, anime.Genres, anime.Score, anime.Popularity,
		)
		if err != nil {
			return fmt.Errorf("failed to upsert anime: %w", err)
		}
	}

	return tx.Commit()
} 