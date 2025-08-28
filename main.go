package main

import (
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func main() {
	const rootPath = "."
	const port = "8080"
	var apiCfg = apiConfig{
		fileserverHits: atomic.Int32{},
	}

	mux := http.NewServeMux()
	handler := http.FileServer(http.Dir(rootPath))
	mux.Handle("/app/", http.StripPrefix("/app", apiCfg.middlewareMetricsInc(handler)))

	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("POST /api/validate_chirp", handlerValidate)

	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerVisiterCount)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerResetVisiterCount)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
