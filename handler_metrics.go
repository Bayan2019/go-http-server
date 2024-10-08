package main

import (
	"fmt"
	"net/http"
)

// write a new middleware method on a *apiConfig 
// that increments the fileserverHits counter every time it's called.
func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}
// Create a new handler that writes the number of requests 
// that have been counted as plain text in this format to the HTTP response
// This handler should be a method on the *apiConfig struct 
// so that it can access the fileserverHits data.
func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	// w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	// Make sure you use the Content-Type header to set the response type to text/html 
	// so that the browser knows how to render it.
	w.Header().Add("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	// w.Write([]byte(fmt.Sprintf("Hits: %d", cfg.fileserverHits.Load())))
	// Swap out the GET /api/metrics endpoint, which just returns plain text, 
	// for a GET /admin/metrics that returns HTML to be rendered in the browser
	w.Write([]byte(fmt.Sprintf(`
		<html>
			<body>
				<h1>Welcome, Chirpy Admin</h1>
				<p>Chirpy has been visited %d times!</p>
			</body>
		</html>
	`, cfg.fileserverHits.Load())))
}