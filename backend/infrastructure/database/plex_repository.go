package database

import (
	"database/sql"
	"time"

	"anime-watchlist/backend/domain"
)

type PlexRepository struct {
	db *sql.DB
}

func NewPlexRepository(db *sql.DB) *PlexRepository {
	return &PlexRepository{db: db}
}

func (r *PlexRepository) UpsertPlexShow(show *domain.PlexShow) error {
	query := `
		INSERT INTO plex_shows (plex_id, title, anilist_id, year, episode_count, last_updated)
		VALUES (?, ?, ?, ?, ?, ?)
		ON CONFLICT(plex_id) DO UPDATE SET
			title = excluded.title,
			anilist_id = excluded.anilist_id,
			year = excluded.year,
			episode_count = excluded.episode_count,
			last_updated = excluded.last_updated
	`

	var anilistID *int
	if show.AnilistID != nil {
		anilistID = show.AnilistID
	}

	_, err := r.db.Exec(query, show.PlexID, show.Title, anilistID, show.Year, show.EpisodeCount, show.LastUpdated)
	return err
}

func (r *PlexRepository) GetPlexShowByPlexID(plexID int) (*domain.PlexShow, error) {
	query := `
		SELECT id, plex_id, title, anilist_id, year, episode_count, last_updated
		FROM plex_shows
		WHERE plex_id = ?
	`

	show := &domain.PlexShow{}
	var anilistID *int

	err := r.db.QueryRow(query, plexID).Scan(
		&show.ID,
		&show.PlexID,
		&show.Title,
		&anilistID,
		&show.Year,
		&show.EpisodeCount,
		&show.LastUpdated,
	)

	if err != nil {
		return nil, err
	}

	show.AnilistID = anilistID
	return show, nil
}

func (r *PlexRepository) GetAllPlexShows() ([]domain.PlexShow, error) {
	query := `
		SELECT id, plex_id, title, anilist_id, year, episode_count, last_updated
		FROM plex_shows
		ORDER BY title
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var shows []domain.PlexShow
	for rows.Next() {
		show := domain.PlexShow{}
		var anilistID *int

		err := rows.Scan(
			&show.ID,
			&show.PlexID,
			&show.Title,
			&anilistID,
			&show.Year,
			&show.EpisodeCount,
			&show.LastUpdated,
		)
		if err != nil {
			return nil, err
		}

		show.AnilistID = anilistID
		shows = append(shows, show)
	}

	return shows, nil
}

func (r *PlexRepository) GetUnmappedShows() ([]domain.PlexShow, error) {
	query := `
		SELECT id, plex_id, title, anilist_id, year, episode_count, last_updated
		FROM plex_shows
		WHERE anilist_id IS NULL
		ORDER BY title
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var shows []domain.PlexShow
	for rows.Next() {
		show := domain.PlexShow{}
		var anilistID *int

		err := rows.Scan(
			&show.ID,
			&show.PlexID,
			&show.Title,
			&anilistID,
			&show.Year,
			&show.EpisodeCount,
			&show.LastUpdated,
		)
		if err != nil {
			return nil, err
		}

		show.AnilistID = anilistID
		shows = append(shows, show)
	}

	return shows, nil
}

func (r *PlexRepository) GetShowsOnServer() (int, error) {
	query := `SELECT COUNT(*) FROM plex_shows`
	
	var count int
	err := r.db.QueryRow(query).Scan(&count)
	return count, err
}

func (r *PlexRepository) GetMappedShowsCount() (int, error) {
	query := `SELECT COUNT(*) FROM plex_shows WHERE anilist_id IS NOT NULL`
	
	var count int
	err := r.db.QueryRow(query).Scan(&count)
	return count, err
}

func (r *PlexRepository) GetUnmappedShowsCount() (int, error) {
	query := `SELECT COUNT(*) FROM plex_shows WHERE anilist_id IS NULL`
	
	var count int
	err := r.db.QueryRow(query).Scan(&count)
	return count, err
}

func (r *PlexRepository) SearchShowsOnServer(searchTerm string) ([]domain.PlexShow, error) {
	query := `
		SELECT id, plex_id, title, anilist_id, year, episode_count, last_updated
		FROM plex_shows
		WHERE title LIKE ?
		ORDER BY title
	`

	rows, err := r.db.Query(query, "%"+searchTerm+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var shows []domain.PlexShow
	for rows.Next() {
		show := domain.PlexShow{}
		var anilistID *int

		err := rows.Scan(
			&show.ID,
			&show.PlexID,
			&show.Title,
			&anilistID,
			&show.Year,
			&show.EpisodeCount,
			&show.LastUpdated,
		)
		if err != nil {
			return nil, err
		}

		show.AnilistID = anilistID
		shows = append(shows, show)
	}

	return shows, nil
}

func (r *PlexRepository) UpdateShowMapping(plexID int, anilistID int) error {
	query := `
		UPDATE plex_shows 
		SET anilist_id = ?, last_updated = ?
		WHERE plex_id = ?
	`

	_, err := r.db.Exec(query, anilistID, time.Now(), plexID)
	return err
} 