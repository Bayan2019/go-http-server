package main

import (
	"net/http"

	"github.com/Bayan2019/go-http-server/internal/auth"
)

// 6. Authentication / 11. Refresh Tokens
func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	// This new endpoint does not accept a request body,
	// but does require a refresh token to be present in the headers,
	// in the same Authorization: Bearer <token> format.
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		//  If it doesn't exist, or if it's expired, respond with a 401 status code.
		respondWithError(w, http.StatusBadRequest, "Couldn't find Refresh Token", err)
		return
	}
	// Revoke the token in the database that matches the token
	// that was passed in the header of the request
	// by setting the revoked_at to the current timestamp.
	// Remember that any time you update a record,
	// you should also be updating the updated_at timestamp.
	_, err = cfg.DB.RevokeToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't Revoke Token", err)
		return
	}

	// no body is returned
	// Respifond with a 204 status code.
	w.WriteHeader(http.StatusNoContent)
}
