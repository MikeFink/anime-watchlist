package application

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"anime-watchlist/backend/domain"
)

type AnilistService struct {
	client  *http.Client
	baseURL string
}

func NewAnilistService() *AnilistService {
	return &AnilistService{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		baseURL: "https://graphql.anilist.co",
	}
}

const searchQuery = `
query ($page: Int, $perPage: Int, $search: String, $status: MediaStatus, $season: MediaSeason, $seasonYear: Int) {
  Page(page: $page, perPage: $perPage) {
    media(type: ANIME, search: $search, status: $status, season: $season, seasonYear: $seasonYear, sort: POPULARITY_DESC, isAdult: false) {
      id
      title {
        romaji
        english
      }
      description
      coverImage {
        large
      }
      bannerImage
      status
      format
      episodes
      duration
      season
      seasonYear
      genres
      averageScore
      popularity
    }
  }
}
`

const animeQuery = `
query ($id: Int) {
  Media(id: $id, type: ANIME) {
    id
    title {
      romaji
      english
    }
    description
    coverImage {
      large
    }
    bannerImage
    status
    format
    episodes
    duration
    season
    seasonYear
    genres
    averageScore
    popularity
  }
}
`

type GraphQLRequest struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables"`
}

func (s *AnilistService) SearchAnime(filter domain.AnimeSearchFilter) ([]domain.Anime, error) {
	variables := map[string]interface{}{
		"page":     filter.Page,
		"perPage":  filter.PageSize,
	}

	if filter.Search != "" {
		variables["search"] = filter.Search
	}
	if filter.Status != "" {
		variables["status"] = filter.Status
	}
	if filter.Season != "" {
		variables["season"] = filter.Season
	}
	if filter.Year > 0 {
		variables["seasonYear"] = filter.Year
	}

	req := GraphQLRequest{
		Query:     searchQuery,
		Variables: variables,
	}

	var response domain.AnilistResponse
	if err := s.makeRequest(req, &response); err != nil {
		return nil, fmt.Errorf("failed to search anime: %w", err)
	}

	animes := make([]domain.Anime, len(response.Data.Page.Media))
	for i, media := range response.Data.Page.Media {
		animes[i] = media.ToDomain()
	}

	return animes, nil
}

func (s *AnilistService) GetAnimeByID(anilistID int) (*domain.Anime, error) {
	variables := map[string]interface{}{
		"id": anilistID,
	}

	req := GraphQLRequest{
		Query:     animeQuery,
		Variables: variables,
	}

	var response struct {
		Data struct {
			Media domain.AnilistAnime `json:"Media"`
		} `json:"data"`
	}

	if err := s.makeRequest(req, &response); err != nil {
		return nil, fmt.Errorf("failed to get anime by ID: %w", err)
	}

	anime := response.Data.Media.ToDomain()
	return &anime, nil
}

func (s *AnilistService) makeRequest(req GraphQLRequest, response interface{}) error {
	jsonData, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequest("POST", s.baseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")

	resp, err := s.client.Do(httpReq)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("anilist API returned status %d: %s", resp.StatusCode, string(body))
	}

	if err := json.Unmarshal(body, response); err != nil {
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return nil
} 