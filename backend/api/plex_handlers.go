package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"anime-watchlist/backend/application"
	"anime-watchlist/backend/domain"
	"anime-watchlist/backend/infrastructure/database"
)

type PlexHandlers struct {
	plexService *application.PlexService
	plexRepo    *database.PlexRepository
}

func NewPlexHandlers(plexService *application.PlexService, plexRepo *database.PlexRepository) *PlexHandlers {
	return &PlexHandlers{
		plexService: plexService,
		plexRepo:    plexRepo,
	}
}

func (h *PlexHandlers) GetServerStatus(w http.ResponseWriter, r *http.Request) {
	status, err := h.plexService.GetServerStatus(h.plexRepo)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get server status", err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, status)
}

func (h *PlexHandlers) SyncPlexShows(w http.ResponseWriter, r *http.Request) {
	shows, err := h.plexService.FetchShowsFromPlex()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to sync plex shows", err.Error())
		return
	}

	for _, show := range shows {
		if err := h.plexRepo.UpsertPlexShow(&show); err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed to save plex show", err.Error())
			return
		}
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Plex shows synced successfully",
		"count":   len(shows),
	})
}

func (h *PlexHandlers) GetShowsOnServer(w http.ResponseWriter, r *http.Request) {
	shows, err := h.plexRepo.GetAllPlexShows()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get shows on server", err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, shows)
}

func (h *PlexHandlers) GetUnmappedShows(w http.ResponseWriter, r *http.Request) {
	shows, err := h.plexRepo.GetUnmappedShows()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get unmapped shows", err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, shows)
}

func (h *PlexHandlers) SearchShowsOnServer(w http.ResponseWriter, r *http.Request) {
	searchTerm := r.URL.Query().Get("q")
	if searchTerm == "" {
		respondWithError(w, http.StatusBadRequest, "Search term is required", "")
		return
	}

	shows, err := h.plexRepo.SearchShowsOnServer(searchTerm)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to search shows", err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, shows)
}

func (h *PlexHandlers) MapShowToAnilist(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed", "")
		return
	}

	var req struct {
		PlexID    int `json:"plex_id"`
		AnilistID int `json:"anilist_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	if err := h.plexRepo.UpdateShowMapping(req.PlexID, req.AnilistID); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to map show", err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Show mapped successfully",
	})
}

func (h *PlexHandlers) AutoMapShow(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed", "")
		return
	}

	var req struct {
		PlexID int `json:"plex_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	show, err := h.plexRepo.GetPlexShowByPlexID(req.PlexID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Show not found", err.Error())
		return
	}

	if err := h.plexService.MapShowToAnilist(show); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to map show", err.Error())
		return
	}

	if show.AnilistID != nil {
		if err := h.plexRepo.UpsertPlexShow(show); err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed to save mapping", err.Error())
			return
		}
	}

	respondWithJSON(w, http.StatusOK, show)
}

func (h *PlexHandlers) BulkAutoMapShows(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed", "")
		return
	}

	var req struct {
		Limit int `json:"limit"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		req.Limit = 50 // Default limit
	}

	if req.Limit <= 0 || req.Limit > 100 {
		req.Limit = 50
	}

	shows, err := h.plexRepo.GetUnmappedShows()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get unmapped shows", err.Error())
		return
	}

	if len(shows) > req.Limit {
		shows = shows[:req.Limit]
	}

	// Convert to pointers for the service
	showPointers := make([]*domain.PlexShow, len(shows))
	for i := range shows {
		showPointers[i] = &shows[i]
	}

	mappedCount, err := h.plexService.BulkAutoMapShows(showPointers)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to bulk map shows", err.Error())
		return
	}

	// Save all successfully mapped shows
	for _, show := range showPointers {
		if show.AnilistID != nil {
			if err := h.plexRepo.UpsertPlexShow(show); err != nil {
				// Log error but continue
				continue
			}
		}
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"message":      "Bulk auto-map completed",
		"processed":    len(shows),
		"mapped":       mappedCount,
		"failed":       len(shows) - mappedCount,
	})
}

func (h *PlexHandlers) CheckShowOnServer(w http.ResponseWriter, r *http.Request) {
	anilistIDStr := r.URL.Query().Get("anilist_id")
	if anilistIDStr == "" {
		respondWithError(w, http.StatusBadRequest, "Anilist ID is required", "")
		return
	}

	anilistID, err := strconv.Atoi(anilistIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid anilist ID", err.Error())
		return
	}

	shows, err := h.plexRepo.GetAllPlexShows()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get shows", err.Error())
		return
	}

	var foundShow *domain.PlexShow
	for _, show := range shows {
		if show.AnilistID != nil && *show.AnilistID == anilistID {
			foundShow = &show
			break
		}
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"on_server": foundShow != nil,
		"show":      foundShow,
	})
} 