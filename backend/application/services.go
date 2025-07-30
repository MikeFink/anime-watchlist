package application

import (
	"fmt"

	"anime-watchlist/backend/domain"
	"anime-watchlist/backend/infrastructure/database"
)

type AnimeService struct {
	watchlistRepo *database.WatchlistRepository
	anilistService *AnilistService
}

func NewAnimeService(watchlistRepo *database.WatchlistRepository, anilistService *AnilistService) *AnimeService {
	return &AnimeService{
		watchlistRepo: watchlistRepo,
		anilistService: anilistService,
	}
}

func (s *AnimeService) SearchAnime(filter domain.AnimeSearchFilter) ([]domain.Anime, error) {
	if err := filter.Validate(); err != nil {
		return nil, err
	}

	animes, err := s.anilistService.SearchAnime(filter)
	if err != nil {
		return nil, fmt.Errorf("failed to search anime: %w", err)
	}

	for i := range animes {
		isWatching, err := s.watchlistRepo.IsInWatchlist(animes[i].AnilistID)
		if err != nil {
			return nil, fmt.Errorf("failed to check watchlist status: %w", err)
		}
		animes[i].IsWatching = isWatching
	}

	return animes, nil
}

func (s *AnimeService) GetWatchlist() ([]domain.Anime, error) {
	watchlistItems, err := s.watchlistRepo.GetWatchlist()
	if err != nil {
		return nil, fmt.Errorf("failed to get watchlist: %w", err)
	}

	animes := make([]domain.Anime, len(watchlistItems))
	for i, item := range watchlistItems {
		anime, err := s.anilistService.GetAnimeByID(item.AnilistID)
		if err != nil {
			return nil, fmt.Errorf("failed to get anime data for ID %d: %w", item.AnilistID, err)
		}
		
		anime.IsWatching = true
		animes[i] = *anime
	}

	return animes, nil
}

func (s *AnimeService) AddToWatchlist(anilistID int) error {
	if anilistID <= 0 {
		return fmt.Errorf("invalid anilist ID: %d", anilistID)
	}

	isWatching, err := s.watchlistRepo.IsInWatchlist(anilistID)
	if err != nil {
		return fmt.Errorf("failed to check watchlist: %w", err)
	}

	if isWatching {
		return fmt.Errorf("anime already in watchlist")
	}

	return s.watchlistRepo.AddToWatchlist(anilistID)
}

func (s *AnimeService) RemoveFromWatchlist(anilistID int) error {
	if anilistID <= 0 {
		return fmt.Errorf("invalid anilist ID: %d", anilistID)
	}

	return s.watchlistRepo.RemoveFromWatchlist(anilistID)
}

func (s *AnimeService) GetWatchlistCount() (int, error) {
	return s.watchlistRepo.GetWatchlistCount()
} 