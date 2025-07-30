package application

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"anime-watchlist/backend/domain"
)

type PlexService struct {
	config domain.PlexConfig
	lastRequest time.Time
}

type PlexShowResponse struct {
	MediaContainer struct {
		Metadata []PlexShowMetadata `json:"Metadata"`
	} `json:"MediaContainer"`
}

type PlexShowMetadata struct {
	RatingKey string `json:"ratingKey"`
	Title     string `json:"title"`
	Year      int    `json:"year"`
	ChildCount int   `json:"childCount"`
}

func NewPlexService(config domain.PlexConfig) *PlexService {
	return &PlexService{
		config: config,
		lastRequest: time.Now().Add(-time.Second), // Allow immediate first request
	}
}

func (s *PlexService) FetchShowsFromPlex() ([]domain.PlexShow, error) {
	if !s.config.SyncEnabled {
		return nil, fmt.Errorf("plex sync is disabled")
	}

	plexURL := fmt.Sprintf("%s/library/sections/%d/all", s.config.ServerURL, s.config.LibraryID)
	
	req, err := http.NewRequest("GET", plexURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("X-Plex-Token", s.config.Token)
	req.Header.Set("Accept", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch from plex: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("plex API returned status: %d", resp.StatusCode)
	}

	var plexResp PlexShowResponse
	if err := json.NewDecoder(resp.Body).Decode(&plexResp); err != nil {
		return nil, fmt.Errorf("failed to decode plex response: %w", err)
	}

	var shows []domain.PlexShow
	for _, metadata := range plexResp.MediaContainer.Metadata {
		plexID, _ := strconv.Atoi(metadata.RatingKey)
		shows = append(shows, domain.PlexShow{
			PlexID:       plexID,
			Title:        metadata.Title,
			Year:         metadata.Year,
			EpisodeCount: metadata.ChildCount,
			LastUpdated:  time.Now(),
		})
	}

	return shows, nil
}

func (s *PlexService) SearchAnilistForShow(title string, year int) (*domain.Anime, error) {
	searchStrategies := []struct {
		searchTerm string
		useYear    bool
		weight     float64
	}{
		{title, true, 1.0},                    // Exact title + year (highest confidence)
		{title, false, 0.8},                   // Exact title without year
		{s.cleanTitle(title), true, 0.9},      // Cleaned title + year
		{s.cleanTitle(title), false, 0.7},     // Cleaned title without year
		{s.extractMainTitle(title), true, 0.6}, // Main title + year
		{s.extractMainTitle(title), false, 0.5}, // Main title without year
	}

	var bestMatch *domain.Anime
	var bestConfidence float64
	const minConfidence = 0.5

	for _, strategy := range searchStrategies {
		anime, err := s.searchAnilistWithStrategy(strategy.searchTerm, strategy.useYear, year)
		if err != nil {
			continue
		}

		if anime != nil {
			confidence := s.calculateConfidence(title, year, anime, strategy.weight)
			if confidence > bestConfidence && confidence >= minConfidence {
				bestMatch = anime
				bestConfidence = confidence
			}
		}
	}

	return bestMatch, nil
}

func (s *PlexService) searchAnilistWithStrategy(searchTerm string, useYear bool, year int) (*domain.Anime, error) {
	// Rate limiting: wait at least 1 second between requests
	timeSinceLast := time.Since(s.lastRequest)
	if timeSinceLast < time.Second {
		time.Sleep(time.Second - timeSinceLast)
	}
	s.lastRequest = time.Now()

	var query string
	var variables map[string]interface{}

	if useYear {
		query = `query SearchAnime($search: String, $year: Int) {
			Page(page: 1, perPage: 10) {
				media(type: ANIME, search: $search, seasonYear: $year, sort: POPULARITY_DESC, isAdult: false) {
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
		}`
		variables = map[string]interface{}{
			"search": searchTerm,
			"year":   year,
		}
	} else {
		query = `query SearchAnime($search: String) {
			Page(page: 1, perPage: 10) {
				media(type: ANIME, search: $search, sort: POPULARITY_DESC, isAdult: false) {
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
		}`
		variables = map[string]interface{}{
			"search": searchTerm,
		}
	}

	req := struct {
		Query     string                 `json:"query"`
		Variables map[string]interface{} `json:"variables"`
	}{
		Query:     query,
		Variables: variables,
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequest("POST", "https://graphql.anilist.co", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create anilist request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to search anilist: %w", err)
	}
	defer resp.Body.Close()

	// Check for rate limiting
	if resp.StatusCode == 429 {
		return nil, fmt.Errorf("rate limited by anilist API")
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("anilist API returned status: %d - %s", resp.StatusCode, string(body))
	}

	var anilistResp domain.AnilistResponse
	if err := json.NewDecoder(resp.Body).Decode(&anilistResp); err != nil {
		return nil, fmt.Errorf("failed to decode anilist response: %w", err)
	}

	if len(anilistResp.Data.Page.Media) == 0 {
		return nil, nil
	}

	anime := anilistResp.Data.Page.Media[0].ToDomain()
	return &anime, nil
}

func (s *PlexService) cleanTitle(title string) string {
	// Remove common prefixes/suffixes and clean up the title
	title = strings.TrimSpace(title)
	
	// Remove quotes
	title = strings.Trim(title, `"'`)
	
	// Remove common movie prefixes
	prefixes := []string{"Eiga ", "Gekijouban ", "Movie: ", "The Movie"}
	for _, prefix := range prefixes {
		if strings.HasPrefix(title, prefix) {
			title = strings.TrimPrefix(title, prefix)
			title = strings.TrimSpace(title)
		}
	}
	
	// Remove common suffixes
	suffixes := []string{" (Movie)", " Movie", " - Movie"}
	for _, suffix := range suffixes {
		if strings.HasSuffix(title, suffix) {
			title = strings.TrimSuffix(title, suffix)
			title = strings.TrimSpace(title)
		}
	}
	
	return title
}

func (s *PlexService) extractMainTitle(title string) string {
	// Extract the main title before any colons, dashes, or parentheses
	title = s.cleanTitle(title)
	
	// Split by common separators and take the first part
	separators := []string{": ", " - ", " (", " ["}
	for _, sep := range separators {
		if idx := strings.Index(title, sep); idx != -1 {
			title = title[:idx]
			title = strings.TrimSpace(title)
		}
	}
	
	return title
}

func (s *PlexService) calculateConfidence(plexTitle string, plexYear int, anime *domain.Anime, baseWeight float64) float64 {
	confidence := baseWeight
	
	// Title similarity
	plexClean := strings.ToLower(s.cleanTitle(plexTitle))
	animeEnglish := strings.ToLower(anime.TitleEnglish)
	animeRomaji := strings.ToLower(anime.TitleRomaji)
	
	// Check exact matches
	if plexClean == animeEnglish || plexClean == animeRomaji {
		confidence += 0.3
	}
	
	// Check if one contains the other
	if strings.Contains(animeEnglish, plexClean) || strings.Contains(animeRomaji, plexClean) ||
	   strings.Contains(plexClean, animeEnglish) || strings.Contains(plexClean, animeRomaji) {
		confidence += 0.2
	}
	
	// Year matching
	if anime.SeasonYear == plexYear {
		confidence += 0.2
	} else if anime.SeasonYear == plexYear+1 || anime.SeasonYear == plexYear-1 {
		confidence += 0.1
	}
	
	// Format matching (TV series vs Movie)
	if strings.Contains(strings.ToLower(plexTitle), "movie") && anime.Format == "MOVIE" {
		confidence += 0.1
	} else if !strings.Contains(strings.ToLower(plexTitle), "movie") && anime.Format == "TV" {
		confidence += 0.1
	}
	
	// Cap confidence at 1.0
	if confidence > 1.0 {
		confidence = 1.0
	}
	
	return confidence
}

func (s *PlexService) GetServerStatus(repo interface {
	GetShowsOnServer() (int, error)
	GetMappedShowsCount() (int, error)
	GetUnmappedShowsCount() (int, error)
}) (*domain.ServerStatus, error) {
	totalShows, err := repo.GetShowsOnServer()
	if err != nil {
		return nil, err
	}

	mappedShows, err := repo.GetMappedShowsCount()
	if err != nil {
		return nil, err
	}

	unmappedShows, err := repo.GetUnmappedShowsCount()
	if err != nil {
		return nil, err
	}

	status := &domain.ServerStatus{
		ShowsOnServer:    totalShows,
		MappedToAnilist:  mappedShows,
		UnmappedShows:    unmappedShows,
		WatchlistShows:   0, // This would need to be calculated separately
		MissingFromServer: 0, // This would need to be calculated separately
	}

	return status, nil
}

func (s *PlexService) MapShowToAnilist(plexShow *domain.PlexShow) error {
	anime, err := s.SearchAnilistForShow(plexShow.Title, plexShow.Year)
	if err != nil {
		return fmt.Errorf("failed to search anilist for %s: %w", plexShow.Title, err)
	}

	if anime != nil {
		plexShow.AnilistID = &anime.AnilistID
		plexShow.Anime = anime
	}

	return nil
}

func (s *PlexService) BulkAutoMapShows(shows []*domain.PlexShow) (int, error) {
	mappedCount := 0
	
	for i, show := range shows {
		if show.AnilistID != nil {
			continue // Already mapped
		}
		
		if err := s.MapShowToAnilist(show); err != nil {
			// If we hit rate limiting, stop processing
			if strings.Contains(err.Error(), "rate limited") {
				break
			}
			continue // Skip on other errors, continue with next show
		}
		
		if show.AnilistID != nil {
			mappedCount++
		}
		
		// Add a small delay every 5 shows to be respectful to the API
		if (i+1)%5 == 0 {
			time.Sleep(2 * time.Second)
		}
	}
	
	return mappedCount, nil
} 