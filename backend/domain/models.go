package domain

import (
	"time"
)

type Anime struct {
	ID            int       `json:"id" db:"id"`
	AnilistID     int       `json:"anilist_id" db:"anilist_id"`
	Title         string    `json:"title" db:"title"`
	TitleEnglish  *string   `json:"title_english" db:"title_english"`
	TitleRomaji   *string   `json:"title_romaji" db:"title_romaji"`
	Description   *string   `json:"description" db:"description"`
	CoverImage    *string   `json:"cover_image" db:"cover_image"`
	BannerImage   *string   `json:"banner_image" db:"banner_image"`
	Status        *string   `json:"status" db:"status"`
	Format        *string   `json:"format" db:"format"`
	Episodes      *int      `json:"episodes" db:"episodes"`
	Duration      *int      `json:"duration" db:"duration"`
	Season        *string   `json:"season" db:"season"`
	SeasonYear    *int      `json:"season_year" db:"season_year"`
	Genres        *string   `json:"genres" db:"genres"`
	Score         *float64  `json:"score" db:"score"`
	Popularity    *int      `json:"popularity" db:"popularity"`
	IsWatching    bool      `json:"is_watching"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}

type WatchlistItem struct {
	ID      int       `json:"id" db:"id"`
	AnimeID int       `json:"anime_id" db:"anime_id"`
	AddedAt time.Time `json:"added_at" db:"added_at"`
	Anime   Anime     `json:"anime"`
}

type AnimeFilter struct {
	Search   string `json:"search"`
	Status   string `json:"status"`
	Season   string `json:"season"`
	Year     int    `json:"year"`
	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
}

func (f *AnimeFilter) Validate() error {
	if f.Page < 1 {
		f.Page = 1
	}
	if f.PageSize < 1 || f.PageSize > 100 {
		f.PageSize = 20
	}
	return nil
}

func (a *Anime) Validate() error {
	if a.Title == "" {
		return ErrInvalidTitle
	}
	if a.AnilistID <= 0 {
		return ErrInvalidAnilistID
	}
	return nil
} 