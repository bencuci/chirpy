package main

import "net/http"

func handlerReadiness(rw http.ResponseWriter, req *http.Request) {
	req.Header.Set("Content-type", "text/plain; charset=utf-8")
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("Checked health: OK"))
}
