package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Bayan2019/go-http-server/internal/auth"
	"github.com/google/uuid"
)

// 8. Webhooks / 1. Webhooks
func (apiCfg *apiConfig) handlerPolkaWebhookRedChirpy(w http.ResponseWriter, r *http.Request) {

	// 8. Webhooks / 4. API Keys
	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find api key", err)
	}
	// It should ensure that the API key in the header matches the one stored in the .env file.
	if apiKey != apiCfg.polkaKey {
		// If it doesn't, the endpoint should respond with a 401 status code.
		respondWithError(w, http.StatusUnauthorized, "API key is invalid", nil)
		return
	}

	// It should accept a request of this shape:
	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID uuid.UUID `json:"user_id"`
		}
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	if params.Event != "user.upgraded" {
		// If the event is anything other than user.upgraded,
		// the endpoint should immediately respond with a 204 status code
		// - we don't care about any other events.
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// If the event is user.upgraded,
	// then it should update the user in the database,
	// and mark that they are a Chirpy Red member.
	_, err = apiCfg.DB.UpdateUserRed(r.Context(), params.Data.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// If the user can't be found, the endpoint should respond with a 404 status code.
			respondWithError(w, http.StatusNotFound, "Couldn't find user", err)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Couldn't update user", err)
		return
	}

	// If the user is upgraded successfully,
	// the endpoint should respond with a 204 status code and an empty response body.
	w.WriteHeader(http.StatusNoContent)
}
