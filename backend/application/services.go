package application

import (
	"fmt"
	"os/exec"

	"anime-watchlist/backend/domain"
	"anime-watchlist/backend/infrastructure/database"
)

type AnimeService struct {
	repo *database.AnimeRepository
}

func NewAnimeService(repo *database.AnimeRepository) *AnimeService {
	return &AnimeService{repo: repo}
}

func (s *AnimeService) GetAllAnime(filter domain.AnimeFilter) ([]domain.Anime, error) {
	return s.repo.GetAll(filter)
}

func (s *AnimeService) GetWatchlist() ([]domain.Anime, error) {
	return s.repo.GetWatchlist()
}

func (s *AnimeService) AddToWatchlist(animeID int) error {
	if animeID <= 0 {
		return fmt.Errorf("invalid anime ID: %d", animeID)
	}

	return s.repo.AddToWatchlist(animeID)
}

func (s *AnimeService) RemoveFromWatchlist(animeID int) error {
	if animeID <= 0 {
		return fmt.Errorf("invalid anime ID: %d", animeID)
	}

	return s.repo.RemoveFromWatchlist(animeID)
}

func (s *AnimeService) SyncAnimeData() error {
	cmd := exec.Command("python3", "scripts/fetch_anime.py")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to sync anime data: %w, output: %s", err, string(output))
	}

	return nil
}

func (s *AnimeService) UpsertAnime(anime *domain.Anime) error {
	return s.repo.Upsert(anime)
}

func (s *AnimeService) BulkUpsertAnime(animes []domain.Anime) error {
	return s.repo.BulkUpsert(animes)
} 