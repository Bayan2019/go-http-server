package main

import (
	"net/http"
	"time"

	"github.com/Bayan2019/go-http-server/internal/auth"
)

// 6. Authentication / 11. Refresh Tokens
func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	// This new endpoint does not accept a request body,
	// but does require a refresh token to be present in the headers,
	// in the same Authorization: Bearer <token> format.
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		//  If it doesn't exist, or if it's expired, respond with a 401 status code.
		respondWithError(w, http.StatusUnauthorized, "Couldn't find Refresh Token", err)
		return
	}

	user, err := cfg.DB.GetUserFromRefreshToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't get user for refresh token", err)
		return
	}

	// respond with a 200 code and this shape:
	type response struct {
		Token string `json:"token"`
	}
	// The token field should be a newly created access token for the given user that expires in 1 hour.
	expirationTime := time.Hour
	accessToken, err := auth.MakeJWT(
		user.ID,
		cfg.jwtSecret,
		expirationTime,
	)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create access JWT", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{Token: accessToken})
}
