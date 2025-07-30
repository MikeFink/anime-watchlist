package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"anime-watchlist/backend/api"
	"anime-watchlist/backend/application"
	"anime-watchlist/backend/domain"
	"anime-watchlist/backend/infrastructure/config"
	"anime-watchlist/backend/infrastructure/database"
)

func main() {
	cfg := config.Load()

	db, err := database.New(cfg.Database.Path)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	watchlistRepo := database.NewWatchlistRepository(db)
	plexRepo := database.NewPlexRepository(db.DB)
	anilistService := application.NewAnilistService()
	
	plexConfig := domain.PlexConfig{
		ServerURL:   cfg.Plex.ServerURL,
		Token:       cfg.Plex.Token,
		LibraryID:   cfg.Plex.LibraryID,
		SyncEnabled: cfg.Plex.SyncEnabled,
	}
	plexService := application.NewPlexService(plexConfig)
	
	service := application.NewAnimeService(watchlistRepo, anilistService)
	handlers := api.NewHandlers(service)
	plexHandlers := api.NewPlexHandlers(plexService, plexRepo)

	mux := http.NewServeMux()

	mux.HandleFunc("/api/anime/watching", handlers.GetWatchlist)
	mux.HandleFunc("/api/anime/search", handlers.SearchAnime)
	mux.HandleFunc("/api/anime/", handlers.HandleWatchlist)
	mux.HandleFunc("/api/watchlist/count", handlers.GetWatchlistCount)

	mux.HandleFunc("/api/plex/status", plexHandlers.GetServerStatus)
	mux.HandleFunc("/api/plex/sync", plexHandlers.SyncPlexShows)
	mux.HandleFunc("/api/plex/shows", plexHandlers.GetShowsOnServer)
	mux.HandleFunc("/api/plex/unmapped", plexHandlers.GetUnmappedShows)
	mux.HandleFunc("/api/plex/search", plexHandlers.SearchShowsOnServer)
	mux.HandleFunc("/api/plex/map", plexHandlers.MapShowToAnilist)
	mux.HandleFunc("/api/plex/auto-map", plexHandlers.AutoMapShow)
	mux.HandleFunc("/api/plex/bulk-auto-map", plexHandlers.BulkAutoMapShows)
	mux.HandleFunc("/api/plex/check", plexHandlers.CheckShowOnServer)

	handler := api.LoggingMiddleware()(
		api.RequestIDMiddleware()(
			api.CORSMiddleware(cfg.CORS.AllowedOrigins, cfg.CORS.AllowedMethods, cfg.CORS.AllowedHeaders)(
				api.StaticFileMiddleware("./static")(
					mux,
				),
			),
		),
	)

	server := &http.Server{
		Addr:    cfg.Server.Host + ":" + cfg.Server.Port,
		Handler: handler,
	}

	go func() {
		log.Printf("Server starting on %s:%s", cfg.Server.Host, cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
} 