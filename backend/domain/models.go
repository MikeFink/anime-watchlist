package domain

import (
	"strings"
	"time"
)

type Anime struct {
	ID            int       `json:"id" db:"id"`
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
	Search   string `json:"search"`
	Status   string `json:"status"`
	Season   string `json:"season"`
	Year     int    `json:"year"`
	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
}

type AnilistAnime struct {
	ID          int     `json:"id"`
	Title       struct {
		Romaji  string `json:"romaji"`
		English string `json:"english"`
	} `json:"title"`
	Description string `json:"description"`
	CoverImage  struct {
		Large string `json:"large"`
	} `json:"coverImage"`
	BannerImage string `json:"bannerImage"`
	Status      string `json:"status"`
	Format      string `json:"format"`
	Episodes    int    `json:"episodes"`
	Duration    int    `json:"duration"`
	Season      string `json:"season"`
	SeasonYear  int    `json:"seasonYear"`
	Genres      []string `json:"genres"`
	AverageScore float64 `json:"averageScore"`
	Popularity  int     `json:"popularity"`
}

type AnilistResponse struct {
	Data struct {
		Page struct {
			Media []AnilistAnime `json:"media"`
		} `json:"Page"`
	} `json:"data"`
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

func (a *Anime) Validate() error {
	if a.Title == "" {
		return ErrInvalidTitle
	}
	if a.AnilistID <= 0 {
		return ErrInvalidAnilistID
	}
	return nil
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