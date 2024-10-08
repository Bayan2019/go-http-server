package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	// "time"

	// "github.com/Bayan2019/rss_blog/internal/auth"
	"github.com/Bayan2019/go-http-server/internal/auth"
	"github.com/Bayan2019/go-http-server/internal/database"
	"github.com/google/uuid"
)

func (apiCfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {

	// It accepts a JSON payload with a body field:
	type parameters struct {
		Body string `json:"body"`
		// It is not an authenticated endpoint
		// User uuid.UUID `json:"user_id"`
	}

	// To post a chirp, a user needs to have valid JWT
	// Use your GetBearerToken
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT", err)
		return
	}
	// Use your ValidateJWT
	userID, err := auth.ValidateJWT(token, apiCfg.jwtSecret)
	if err != nil {
		// If the JWT is invalid, return a 401 Unauthorized response
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT", err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}

	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %s", err), err)
		return
	}

	// Delete the /api/validate_chirp endpoint that we created before,
	// but port all that logic into this one.
	// Users should not be allowed to create invalid chirps!
	cleaned, err := validateChirp(params.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error(), err)
		return
	}
	// If the Chirp is valid, respond with a 200 code and this body:
	chirp, err := apiCfg.DB.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   cleaned,
		UserID: userID,
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't create feed : %s", err), err)
		return
	}

	// feedFollow, err := apiCfg.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
	// 	ID:        uuid.New(),
	// 	CreatedAt: time.Now().UTC(),
	// 	UpdatedAt: time.Now().UTC(),
	// 	UserID:    user.ID,
	// 	FeedID:    feed.ID,
	// })
	// if err != nil {
	// 	respondWithError(w, 400, fmt.Sprintf("Couldn't create feed follow: %s", err))
	// 	return
	// }

	respondWithJSON(w, http.StatusCreated, databaseChirpToChirp(chirp))
}

func validateChirp(body string) (string, error) {
	const maxChirpLength = 140
	if len(body) > maxChirpLength {
		// if the Chirp is too long, respond with a 400 code and this body:
		// respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		// return
		return "", errors.New("Chirp is too long")
	}
	// Assuming the length validation passed,
	// replace any of the following words in the Chirp with the static 4-character string ****
	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}
	cleaned := getCleanedBody(body, badWords)
	return cleaned, nil
}

func getCleanedBody(body string, badWords map[string]struct{}) string {
	words := strings.Split(body, " ")
	for i, word := range words {
		// Be sure to match against uppercase versions of the words as well
		loweredWord := strings.ToLower(word)
		if _, ok := badWords[loweredWord]; ok {
			// to replace all "profane" words with 4 asterisks: ****
			words[i] = "****"
		}
	}
	cleaned := strings.Join(words, " ")
	return cleaned
}

// 5. Storage 11. Get All Chirps
// Add a GET /api/chirps endpoint that returns all chirps in the database.
// Order them by created_at in ascending order.
func (apiCfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	// 9. Documentation 4. Sorting Chirps
	//  It should accept an optional query parameter called sort
	sort := r.URL.Query().Get("sort")
	if sort != "desc" {
		// asc is the default if no sort query parameter is provided.
		sort = "asc"
	}
	// 9. Documentation 1. Documentation
	// Update the GET /api/chirps endpoint. It should accept an optional query parameter called author_id.
	authorIDstr := r.URL.Query().Get("author_id")
	if authorIDstr == "" {
		// If the author_id query parameter is not provided,
		// the endpoint should return all chirps as it did before.
		if sort == "asc" {
			chirps, err := apiCfg.DB.GetChirpsAsc(r.Context())
			if err != nil {
				respondWithError(w, 400, fmt.Sprintf("Couldn't get chirps : %s", err), err)
				return
			}

			respondWithJSON(w, http.StatusOK, databaseChirpsToChirps(chirps))
			return
		}

		chirps, err := apiCfg.DB.GetChirpsDesc(r.Context())
		if err != nil {
			respondWithError(w, 400, fmt.Sprintf("Couldn't get chirps : %s", err), err)
			return
		}

		respondWithJSON(w, http.StatusOK, databaseChirpsToChirps(chirps))
		return
	}
	// If the author_id query parameter is provided,
	// the endpoint should return only the chirps for that author.
	authorID, err := uuid.Parse(authorIDstr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid author ID", err)
		return
	}

	if sort == "asc" {
		chirpsByAuthor, err := apiCfg.DB.GetChirpsByAuthorAsc(r.Context(), authorID)
		if err != nil {
			respondWithError(w, 400, fmt.Sprintf("Couldn't get feeds : %s", err), err)
		}

		respondWithJSON(w, http.StatusOK, databaseChirpsToChirps(chirpsByAuthor))
		return
	}

	chirpsByAuthor, err := apiCfg.DB.GetChirpsByAuthorDesc(r.Context(), authorID)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't get feeds : %s", err), err)
	}

	respondWithJSON(w, http.StatusOK, databaseChirpsToChirps(chirpsByAuthor))
}

func (apiCfg *apiConfig) handlerGetChirp(w http.ResponseWriter, r *http.Request) {
	// You can get the string value of the path parameter like in Go
	// with the http.Request.PathValue method.
	chirpIDstr := r.PathValue("chirpID")

	chirpID, err := uuid.Parse(chirpIDstr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID", err)
		return
	}

	dbChirp, err := apiCfg.DB.GetChirp(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't get chirp", err)
		return
	}

	// If the chirp is found, return it like so with a 200 code:
	respondWithJSON(w, http.StatusOK, databaseChirpToChirp(dbChirp))
}

// 7. Authorization / 4. Delete Chirp
// Add a new DELETE /api/chirps/{chirpID} route to your server
// that deletes a chirp from the database by its id.
func (apiCfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {
	// To delete a chirp, a user needs to have valid JWT
	// This is an authenticated endpoint,
	// so be sure to check the token in the header.
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT", err)
		return
	}
	// Use your ValidateJWT
	userID, err := auth.ValidateJWT(token, apiCfg.jwtSecret)
	if err != nil {
		// If the JWT is invalid, return a 401 Unauthorized response
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT", err)
		return
	}

	// You can get the string value of the path parameter like in Go
	// with the http.Request.PathValue method.
	chirpIDstr := r.PathValue("chirpID")

	chirpID, err := uuid.Parse(chirpIDstr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID", err)
		return
	}

	dbChirp, err := apiCfg.DB.GetChirp(r.Context(), chirpID)
	if err != nil {
		// If the chirp is not found, return a 404 status code.
		respondWithError(w, http.StatusNotFound, "Couldn't get chirp", err)
		return
	}

	// Only allow the deletion of a chirp
	// if the user is the author of the chirp.
	if userID != dbChirp.UserID {
		// If they are not, return a 403 status code.
		respondWithError(w, http.StatusForbidden, "Not an author of the chirp", nil)
		return
	}
	err = apiCfg.DB.DeleteChirp(r.Context(), dbChirp.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Chirp is not deleted", err)
		return
	}
	// If the chirp is deleted successfully, return a 204 status code.
	w.WriteHeader(http.StatusNoContent)
}
