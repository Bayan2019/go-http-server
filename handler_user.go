package main

import (
	"encoding/json"
	"net/http"

	"github.com/Bayan2019/go-http-server/internal/auth"
	"github.com/Bayan2019/go-http-server/internal/database"
)

func (apiCfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {

	// It accepts an email as JSON in the request body
	type parameters struct {
		Email string `json:"email"`
		// The body parameters should now require a new password field:
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}

	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, "Error parsing JSON", err)
		return
	}

	// 6. Authentication / 1. Authentication with Passwords
	// Hash the password using the bcrypt.GenerateFromPassword function
	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't hash password", err)
		return
	}

	user, err := apiCfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		respondWithError(w, 400, "Couldn't create user", err)
	}

	respondWithJSON(w, 201, databaseUserToUser(user))
}

// 7. Authorization / 1. Authorization
// dd a PUT /api/users endpoint so
// that users can update their own (but not other's) email and password.
func (apiCfg *apiConfig) handlerEditUser(w http.ResponseWriter, r *http.Request) {
	// An access token in the header
	// To change info, a user needs to have valid JWT
	// Use your GetBearerToken
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		// If the access token is malformed or missing,
		// respond with a 401 status code.
		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT", err)
		return
	}
	// Use your ValidateJWT
	userID, err := auth.ValidateJWT(token, apiCfg.jwtSecret)
	if err != nil {
		// If the JWT is invalid, return a 401 Unauthorized response
		// If the access token is malformed or missing,
		// respond with a 401 status code.
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT", err)
		return
	}

	// A new password and email in the request body
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	// Hash the password using the bcrypt.GenerateFromPassword function
	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't hash password", err)
		return
	}

	// update the hashed password and the email
	// for the authenticated user in the database
	user, err := apiCfg.DB.UpdateUser(r.Context(), database.UpdateUserParams{
		ID:             userID,
		Email:          params.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't update user", err)
	}

	// Respond with a 200
	// if everything is successful
	// and the newly updated User resource
	// (omitting the password of course).
	respondWithJSON(w, http.StatusOK, databaseUserToUser(user))
}

// func (apiCfg *apiConfig) handlerGetUser(w http.ResponseWriter, r *http.Request, user database.User) {

// 	respondWithJSON(w, 200, databaseUserToUser(user))
// }
