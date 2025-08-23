package main

import (
	"log"
	"net/http"
)

func main() {
	const rootPath = "."
	const port = "8080"

	mux := http.NewServeMux()
	mux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir(rootPath))))
	mux.HandleFunc("/healthz", handlerReadiness)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

func handlerReadiness(respWriter http.ResponseWriter, req *http.Request) {
	req.Header.Set("Content-type", "text/plain; charset=utf-8")
	respWriter.WriteHeader(http.StatusOK)
	respWriter.Write([]byte("Checked health: OK"))
}
