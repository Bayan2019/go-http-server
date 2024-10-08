package main

import "net/http"

// /reset will need to be a method on the *apiConfig struct
// so that it can also access the fileserverHits
func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	if cfg.Platform != "dev" {
		// If PLATFORM is not equal to "dev",
		// this endpoint should return a 403 Forbidden
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Reset is only allowed in dev environment."))
		return
	}

	cfg.fileserverHits.Store(0)
	// to delete all users in the database
	// (but don't mess with the schema)
	cfg.DB.Reset(r.Context())
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0 and database reset to initial state."))
}
