package domain

import (
	"strings"
	"time"
)

type Anime struct {
	AnilistID     int       `json:"anilist_id" db:"anilist_id"`
	Title         string    `json:"title" db:"title"`
	TitleEnglish  string    `json:"title_english" db:"title_english"`
	TitleRomaji   string    `json:"title_romaji" db:"title_romaji"`
	Description   string    `json:"description" db:"description"`
	CoverImage    string    `json:"cover_image" db:"cover_image"`
	BannerImage   string    `json:"banner_image" db:"banner_image"`
	Status        string    `json:"status" db:"status"`
	Format        string    `json:"format" db:"format"`
	Episodes      int       `json:"episodes" db:"episodes"`
	Duration      int       `json:"duration" db:"duration"`
	Season        string    `json:"season" db:"season"`
	SeasonYear    int       `json:"season_year" db:"season_year"`
	Genres        string    `json:"genres" db:"genres"`
	Score         float64   `json:"score" db:"score"`
	Popularity    int       `json:"popularity" db:"popularity"`
	IsWatching    bool      `json:"is_watching"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}

type WatchlistItem struct {
	ID        int       `json:"id" db:"id"`
	AnilistID int       `json:"anilist_id" db:"anilist_id"`
	AddedAt   time.Time `json:"added_at" db:"added_at"`
	Anime     *Anime    `json:"anime,omitempty"`
}

type AnimeSearchFilter struct {
	Search     string `json:"search"`
	Status     string `json:"status"`
	Season     string `json:"season"`
	SeasonYear int    `json:"season_year"`
	Page       int    `json:"page"`
	PageSize   int    `json:"page_size"`
}

func (f *AnimeSearchFilter) Validate() error {
	if f.Page < 1 {
		f.Page = 1
	}
	if f.PageSize < 1 || f.PageSize > 50 {
		f.PageSize = 20
	}
	return nil
}

type PlexShow struct {
	ID          int       `json:"id" db:"id"`
	PlexID      int       `json:"plex_id" db:"plex_id"`
	Title       string    `json:"title" db:"title"`
	AnilistID   *int      `json:"anilist_id" db:"anilist_id"`
	Year        int       `json:"year" db:"year"`
	EpisodeCount int      `json:"episode_count" db:"episode_count"`
	LastUpdated time.Time `json:"last_updated" db:"last_updated"`
	Anime       *Anime    `json:"anime,omitempty"`
}

type PlexConfig struct {
	ServerURL   string `json:"server_url"`
	Token       string `json:"token"`
	LibraryID   int    `json:"library_id"`
	SyncEnabled bool   `json:"sync_enabled"`
}

type ServerStatus struct {
	ShowsOnServer    int `json:"shows_on_server"`
	MappedToAnilist  int `json:"mapped_to_anilist"`
	UnmappedShows    int `json:"unmapped_shows"`
	WatchlistShows   int `json:"watchlist_shows"`
	MissingFromServer int `json:"missing_from_server"`
}

type AnilistAnime struct {
	ID          int     `json:"id"`
	Title       Title   `json:"title"`
	Description string  `json:"description"`
	CoverImage  Cover   `json:"coverImage"`
	BannerImage string  `json:"bannerImage"`
	Status      string  `json:"status"`
	Format      string  `json:"format"`
	Episodes    int     `json:"episodes"`
	Duration    int     `json:"duration"`
	Season      string  `json:"season"`
	SeasonYear  int     `json:"seasonYear"`
	Genres      []string `json:"genres"`
	AverageScore float64 `json:"averageScore"`
	Popularity   int     `json:"popularity"`
}

type Title struct {
	Romaji  string `json:"romaji"`
	English string `json:"english"`
}

type Cover struct {
	Large string `json:"large"`
}

type AnilistResponse struct {
	Data struct {
		Page struct {
			Media []AnilistAnime `json:"media"`
		} `json:"Page"`
	} `json:"data"`
}

func (a *AnilistAnime) ToDomain() Anime {
	genres := ""
	if len(a.Genres) > 0 {
		genres = strings.Join(a.Genres, ", ")
	}

	return Anime{
		AnilistID:    a.ID,
		Title:        a.Title.Romaji,
		TitleEnglish: a.Title.English,
		TitleRomaji:  a.Title.Romaji,
		Description:  a.Description,
		CoverImage:   a.CoverImage.Large,
		BannerImage:  a.BannerImage,
		Status:       a.Status,
		Format:       a.Format,
		Episodes:     a.Episodes,
		Duration:     a.Duration,
		Season:       a.Season,
		SeasonYear:   a.SeasonYear,
		Genres:       genres,
		Score:        a.AverageScore,
		Popularity:   a.Popularity,
	}
} 