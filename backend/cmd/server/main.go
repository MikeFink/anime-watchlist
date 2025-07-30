package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"anime-watchlist/backend/api"
	"anime-watchlist/backend/application"
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

	repo := database.NewAnimeRepository(db)
	service := application.NewAnimeService(repo)
	handlers := api.NewHandlers(service)

	mux := http.NewServeMux()

	mux.HandleFunc("/api/anime/watching", handlers.GetWatchlist)
	mux.HandleFunc("/api/anime", handlers.GetAnime)
	mux.HandleFunc("/api/anime/", handlers.HandleWatchlist)
	mux.HandleFunc("/api/sync", handlers.SyncAnimeData)

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