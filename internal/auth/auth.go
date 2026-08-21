package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func HashPassword(pass string) (string, error) {
	hashedPass, err := argon2id.CreateHash(pass, argon2id.DefaultParams)
	if err != nil {
		return "", err
	}
	return hashedPass, nil
}

func CheckPasswordHash(pass, hash string) (bool, error) {
	match, err := argon2id.ComparePasswordAndHash(pass, hash)
	return match, err
}

func MakeJWT(userID uuid.UUID, tokenSecret string) (string, error) {
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "music-box-access",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Hour)),
		Subject:   userID.String(),
	})

	result, err := jwtToken.SignedString([]byte(tokenSecret))
	return result, err

}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	claims := &jwt.RegisteredClaims{}

	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(token *jwt.Token) (any, error) {
			return []byte(tokenSecret), nil
		},
	)
	if err != nil {
		return uuid.Nil, err
	}
	if token == nil {
		TError := fmt.Errorf("Token is nil")
		return uuid.Nil, TError
	}
	if !token.Valid {
		TError := fmt.Errorf("Token is not valid")
		return uuid.Nil, TError
	}
	sub, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.Nil, err
	}

	ID, err := uuid.Parse(sub)
	return ID, err
}

func GetBearerToken(headers http.Header) (string, error) {
	auth := headers.Get("Authorization")
	if auth == "" {
		return "", fmt.Errorf("missing authorization header")
	}
	const prefix = "Bearer "
	if !strings.HasPrefix(auth, prefix) {
		return "", fmt.Errorf("malformed authorization header")
	}
	return strings.TrimSpace(strings.TrimPrefix(auth, prefix)), nil
}

func MakeRefreshToken() (string, error) {
	token := make([]byte, 32)
	_, err := rand.Read(token)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(token), nil
}

func GetAPIKey(headers http.Header) (string, error) {
	auth := headers.Get("Authorization")
	if auth == "" {
		TError := fmt.Errorf("A Api Error occurred")
		return "", TError
	}
	return strings.TrimSpace(strings.TrimPrefix(auth, "ApiKey")), nil
}
