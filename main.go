package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/bencuci/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	dbQueries      *database.Queries
	platform       string
	secret         string
}

func main() {
	godotenv.Load()
	const rootPath = "."
	const port = "8080"
	dbURL := os.Getenv("DB_URL")
	platform := os.Getenv("PLATFORM")
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable is not set")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error opening the database: %v", err)
	}

	var apiCfg = apiConfig{
		fileserverHits: atomic.Int32{},
		dbQueries:      database.New(db),
		platform:       platform,
		secret:         jwtSecret,
	}

	mux := http.NewServeMux()
	handler := http.FileServer(http.Dir(rootPath))
	mux.Handle("/app/", http.StripPrefix("/app", apiCfg.middlewareMetricsInc(handler)))

	mux.HandleFunc("GET /api/healthz", handlerReadiness)

	mux.HandleFunc("POST /api/login", apiCfg.handlerLogin)
	mux.HandleFunc("POST /api/refresh", apiCfg.handlerRefreshToken)
	mux.HandleFunc("POST /api/revoke", apiCfg.handlerRevokeRefreshToken)

	mux.HandleFunc("POST /api/users", apiCfg.handlerCreateUser)
	mux.HandleFunc("PUT /api/users", apiCfg.handlerUpdateUser)

	mux.HandleFunc("POST /api/chirps", apiCfg.handlerPostChirp)
	mux.HandleFunc("GET /api/chirps", apiCfg.handlerGetChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.handlerGetChirp)
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", apiCfg.handlerDeleteChirp)

	mux.HandleFunc("POST /admin/reset", apiCfg.handlerResetVisiterCount)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerVisiterCount)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
