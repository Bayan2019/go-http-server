package main

import (
	"time"

	"github.com/Bayan2019/go-http-server/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	// Do NOT return the hashed password in the response
	HashedPassword string `json:"-"`
	// APIKey    string    `json:"api_key"`
	IsChirpyRed bool `json:"is_chirpy_red"`
}

func databaseUserToUser(dbUser database.User) User {
	return User{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Email:     dbUser.Email,
		// HashedPassword: dbUser.HashedPassword,
		// APIKey:    dbUser.ApiKey,
		IsChirpyRed: dbUser.IsChirpyRed,
	}
}

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func databaseChirpToChirp(dbChirp database.Chirp) Chirp {
	return Chirp{
		ID:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body:      dbChirp.Body,
		UserID:    dbChirp.UserID,
	}
}

func databaseChirpsToChirps(dbChirps []database.Chirp) []Chirp {
	chirps := []Chirp{}
	for _, dbChirp := range dbChirps {
		chirps = append(chirps, databaseChirpToChirp(dbChirp))
	}

	return chirps
}

type RefreshToken struct {
	Token     string    `json:"token"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserID    uuid.UUID `json:"user_id"`
	ExpiresAt time.Time `json:"expires_at"`
	RevokedAt time.Time `json:"revoked_at"`
}

func databaseRefreshTokenToRefreshToken(dbRefreshToken database.RefreshToken) RefreshToken {
	return RefreshToken{
		Token:     dbRefreshToken.Token,
		CreatedAt: dbRefreshToken.CreatedAt,
		UpdatedAt: dbRefreshToken.UpdatedAt,
		UserID:    dbRefreshToken.UserID,
		ExpiresAt: dbRefreshToken.ExpiresAt,
		RevokedAt: dbRefreshToken.RevokedAt.Time,
	}
}
