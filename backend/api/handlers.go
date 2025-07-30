package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"anime-watchlist/backend/application"
	"anime-watchlist/backend/domain"
)

type Handlers struct {
	animeService *application.AnimeService
}

func NewHandlers(animeService *application.AnimeService) *Handlers {
	return &Handlers{animeService: animeService}
}

func (h *Handlers) SearchAnime(w http.ResponseWriter, r *http.Request) {
	filter := domain.AnimeSearchFilter{
		Search:   r.URL.Query().Get("search"),
		Status:   r.URL.Query().Get("status"),
		Season:   r.URL.Query().Get("season"),
		Page:     1,
		PageSize: 20,
	}

	if yearStr := r.URL.Query().Get("year"); yearStr != "" {
		if year, err := strconv.Atoi(yearStr); err == nil {
			filter.SeasonYear = year
		}
	}

	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil && page > 0 {
			filter.Page = page
		}
	}

	if pageSizeStr := r.URL.Query().Get("page_size"); pageSizeStr != "" {
		if pageSize, err := strconv.Atoi(pageSizeStr); err == nil && pageSize > 0 {
			filter.PageSize = pageSize
		}
	}

	animes, err := h.animeService.SearchAnime(filter)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to search anime", err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, animes)
}

func (h *Handlers) GetWatchlist(w http.ResponseWriter, r *http.Request) {
	animes, err := h.animeService.GetWatchlist()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to fetch watchlist", err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, animes)
}

func (h *Handlers) HandleWatchlist(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathParts) < 4 {
		respondWithError(w, http.StatusBadRequest, "Invalid path", "Expected /api/anime/{anilist_id}/watch")
		return
	}

	anilistIDStr := pathParts[2]
	anilistID, err := strconv.Atoi(anilistIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid anilist ID", "Anilist ID must be a valid integer")
		return
	}

	switch r.Method {
	case "POST":
		err = h.animeService.AddToWatchlist(anilistID)
		if err != nil {
			status := http.StatusInternalServerError
			message := "Failed to add anime to watchlist"

			if err.Error() == "anime already in watchlist" {
				status = http.StatusConflict
				message = "Anime already in watchlist"
			}

			respondWithError(w, status, message, err.Error())
			return
		}

		respondWithJSON(w, http.StatusOK, map[string]string{"message": "Added to watchlist"})

	case "DELETE":
		err = h.animeService.RemoveFromWatchlist(anilistID)
		if err != nil {
			status := http.StatusInternalServerError
			message := "Failed to remove anime from watchlist"

			if err.Error() == "anime not in watchlist" {
				status = http.StatusNotFound
				message = "Anime not in watchlist"
			}

			respondWithError(w, status, message, err.Error())
			return
		}

		respondWithJSON(w, http.StatusOK, map[string]string{"message": "Removed from watchlist"})

	default:
		respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed", "")
	}
}

func (h *Handlers) GetWatchlistCount(w http.ResponseWriter, r *http.Request) {
	count, err := h.animeService.GetWatchlistCount()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get watchlist count", err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]int{"count": count})
}

func respondWithJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func respondWithError(w http.ResponseWriter, status int, error, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{
		"error":   error,
		"message": message,
	})
} 