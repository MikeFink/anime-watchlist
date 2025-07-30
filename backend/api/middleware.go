package api

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func CORSMiddleware(allowedOrigins, allowedMethods, allowedHeaders []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			if origin != "" {
				for _, allowedOrigin := range allowedOrigins {
					if allowedOrigin == "*" || allowedOrigin == origin {
						w.Header().Set("Access-Control-Allow-Origin", origin)
						break
					}
				}
			}

			if r.Method == "OPTIONS" {
				w.Header().Set("Access-Control-Allow-Methods", joinHeaders(allowedMethods))
				w.Header().Set("Access-Control-Allow-Headers", joinHeaders(allowedHeaders))
				w.Header().Set("Access-Control-Allow-Credentials", "true")
				w.Header().Set("Access-Control-Max-Age", "86400")
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func RequestIDMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get("X-Request-ID")
			if requestID == "" {
				requestID = generateRequestID()
			}
			w.Header().Set("X-Request-ID", requestID)
			next.ServeHTTP(w, r)
		})
	}
}

func LoggingMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			next.ServeHTTP(w, r)
			log.Printf("%s - [%s] \"%s %s %s\" %d %s",
				r.RemoteAddr,
				time.Now().Format(time.RFC1123),
				r.Method,
				r.URL.Path,
				r.Proto,
				http.StatusOK,
				time.Since(start),
			)
		})
	}
}

func StaticFileMiddleware(staticPath string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.HasPrefix(r.URL.Path, "/static/") {
				requestedPath := strings.TrimPrefix(r.URL.Path, "/static/")
				fullPath := filepath.Join(staticPath, requestedPath)

				if _, err := os.Stat(fullPath); os.IsNotExist(err) {
					http.NotFound(w, r)
					return
				}

				content, err := os.ReadFile(fullPath)
				if err != nil {
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					return
				}

				switch {
				case strings.HasSuffix(requestedPath, ".css"):
					w.Header().Set("Content-Type", "text/css")
				case strings.HasSuffix(requestedPath, ".js"):
					w.Header().Set("Content-Type", "application/javascript")
				case strings.HasSuffix(requestedPath, ".png"):
					w.Header().Set("Content-Type", "image/png")
				case strings.HasSuffix(requestedPath, ".jpg"), strings.HasSuffix(requestedPath, ".jpeg"):
					w.Header().Set("Content-Type", "image/jpeg")
				case strings.HasSuffix(requestedPath, ".gif"):
					w.Header().Set("Content-Type", "image/gif")
				case strings.HasSuffix(requestedPath, ".svg"):
					w.Header().Set("Content-Type", "image/svg+xml")
				case strings.HasSuffix(requestedPath, ".ico"):
					w.Header().Set("Content-Type", "image/x-icon")
				case strings.HasSuffix(requestedPath, ".json"):
					w.Header().Set("Content-Type", "application/json")
				case strings.HasSuffix(requestedPath, ".txt"):
					w.Header().Set("Content-Type", "text/plain")
				case strings.HasSuffix(requestedPath, ".map"):
					w.Header().Set("Content-Type", "application/json")
				}

				w.WriteHeader(http.StatusOK)
				w.Write(content)
				return
			}

			if !strings.HasPrefix(r.URL.Path, "/api/") {
				indexPath := filepath.Join(staticPath, "index.html")
				if _, err := os.Stat(indexPath); err == nil {
					http.ServeFile(w, r, indexPath)
					return
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}

func joinHeaders(headers []string) string {
	result := ""
	for i, header := range headers {
		if i > 0 {
			result += ", "
		}
		result += header
	}
	return result
}

func generateRequestID() string {
	bytes := make([]byte, 8)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
} 