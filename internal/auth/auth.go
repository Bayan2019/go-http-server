package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type TokenType string

const (
	// TokenTypeAccess -
	// Set the Issuer to "chirpy"
	TokenTypeAccess TokenType = "chirpy-access"
)

// ErrNoAuthHeaderIncluded -
var ErrNoAuthHeaderIncluded = errors.New("no auth header included in request")

// 6. Authentication / 1. Authentication with Passwords
// Hash the password using the bcrypt.GenerateFromPassword function
// HashPassword -
func HashPassword(password string) (string, error) {
	dat, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(dat), nil
}

// 6. Authentication / 1. Authentication with Passwords
// Use the bcrypt.CompareHashAndPassword function
// to compare the password that the user entered in the HTTP request
// with the password that is stored in the database.
// CheckPasswordHash -
func CheckPasswordHash(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

// 6. Authentication / 1. Authentication with JWTs
// Add a MakeJWT function to your auth package:
// MakeJWT -
func MakeJWT(
	userID uuid.UUID,
	tokenSecret string,
	expiresIn time.Duration,
) (string, error) {
	signingKey := []byte(tokenSecret)
	// Use jwt.NewWithClaims to create a new token
	// Use jwt.SigningMethodHS256 as the signing method.
	// Use jwt.RegisteredClaims as the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer: string(TokenTypeAccess),
		// Set IssuedAt to the current time in UTC
		IssuedAt: jwt.NewNumericDate(time.Now().UTC()),
		// Set ExpiresAt to the current time plus the expiration time (expiresIn)
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
		// Set the Subject to a stringified version of the user's id
		Subject: userID.String(),
	})
	// Use token.SignedString to sign the token with the secret key.
	return token.SignedString(signingKey)
}

// 6. Authentication / 1. Authentication with JWTs
// Add a ValidateJWT function to your auth package:
// ValidateJWT -
func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	claimsStruct := jwt.RegisteredClaims{}
	// Use the jwt.ParseWithClaims function
	// to validate the signature of the JWT and extract the claims
	// into a *jwt.Token struct.
	token, err := jwt.ParseWithClaims(
		tokenString,
		&claimsStruct,
		func(token *jwt.Token) (interface{}, error) { return []byte(tokenSecret), nil },
	)
	if err != nil {
		// An error will be returned if the token is invalid or has expired.
		// If the token is invalid, return a 401 Unauthorized response from your handler.
		return uuid.Nil, err
	}

	// If all is well with the token, use the token.Claims interface
	// to get access to the user's id from the claims
	// (which should be stored in the Subject field).
	userIDString, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.Nil, err
	}

	issuer, err := token.Claims.GetIssuer()
	if err != nil {
		return uuid.Nil, err
	}
	if issuer != string(TokenTypeAccess) {
		return uuid.Nil, errors.New("invalid issuer")
	}

	id, err := uuid.Parse(userIDString)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid user ID: %w", err)
	}
	// Return the id as a uuid.UUID
	return id, nil
}

// 6. Authentication / 1. Authentication with JWTs
// Add a GetBearerToken function to your auth package
// GetBearerToken -
func GetBearerToken(headers http.Header) (string, error) {
	// Auth information will come into our server in the Authorization header.
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		// If the header doesn't exist, return an error.
		return "", ErrNoAuthHeaderIncluded
	}
	// stripping off the Bearer prefix and whitespace
	splitAuth := strings.Split(authHeader, " ")
	if len(splitAuth) < 2 || splitAuth[0] != "Bearer" {
		// If the header doesn't exist, return an error.
		return "", errors.New("malformed authorization header")
	}
	// return the TOKEN_STRING if it exists
	return splitAuth[1], nil
}

// 6. Authentication / 11. Refresh Tokens
// Add a func MakeRefreshToken() (string, error) function to your internal/auth package.
func MakeRefreshToken() (string, error) {
	token := make([]byte, 32)
	// It should use the following to generate a random 256-bit (32-byte) hex-encoded string:
	// rand.Read to generate 32 bytes (256 bits)
	// of random data from the crypto/rand package
	_, err := rand.Read(token)
	if err != nil {
		return "", err
	}
	// hex.EncodeToString to convert the random data to a hex string
	refreshToken := hex.EncodeToString(token)
	return refreshToken, nil
}

// 8. Webhooks / 4. API Keys
// Add a func GetAPIKey(headers http.Header) (string, error) to your auth package.
func GetAPIKey(headers http.Header) (string, error) {
	//  It should extract the api key from the Authorization header
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", ErrNoAuthHeaderIncluded
	}
	// You'll need to strip out the ApiKey part and the whitespace and return just the key.
	splitAuth := strings.Split(authHeader, " ")
	if len(splitAuth) < 2 || splitAuth[0] != "ApiKey" {
		return "", errors.New("malformed authorization header")
	}

	return splitAuth[1], nil
}
