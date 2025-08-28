package main

import "net/http"

func (cfg *apiConfig) handlerResetVisiterCount(rw http.ResponseWriter, req *http.Request) {
	req.Header.Set("Content-type", "text/plain; charset=utf-8")
	cfg.fileserverHits.Store(0)
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("Resetted"))
}
