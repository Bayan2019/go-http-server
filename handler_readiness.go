package main

import "net/http"

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	// Write the Content-Type header
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	// Write the status code using w.WriteHeader
	w.WriteHeader(http.StatusOK)
	// Write the body text using w.Write
	w.Write([]byte(http.StatusText(http.StatusOK)))
}