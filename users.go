package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/ErisLietus/go-http-client/internal/auth"
	"github.com/ErisLietus/go-http-client/internal/database"
	"github.com/google/uuid"
)

type CreatedUser struct {
	ID        uuid.UUID `json:"id"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	Expires   int       `json:"expires_in_seconds"`
}

type FoundUser struct {
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
	IsChirpyRed  bool      `json:"is_chirpy_red"`
}

type HashedUser struct {
	Email          string
	HashedPassword string
}
type UpdateUser struct {
	Email    string
	Password string
}

func (cfg *apiConfig) handlerUsersCreate(w http.ResponseWriter, r *http.Request) {
	var user CreatedUser
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&user); err != nil {
		respondWithError(w, 400, "Something when wrong")
		return
	}

	ctx := r.Context()

	hashed, err := auth.HashPassword(user.Password)
	if err != nil {
		respondWithError(w, 400, "hashing error")
	}

	params := database.CreateUserParams{
		Email:          user.Email,
		HashedPassword: hashed,
	}

	data, err := cfg.db.CreateUser(ctx, params)
	if err != nil {
		log.Printf("CreateUser failed: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Server Error")
		return
	}

	response := FoundUser{
		ID:        data.ID,
		CreatedAt: data.CreatedAt,
		UpdatedAt: data.UpdatedAt,
		Email:     data.Email,
	}
	respondWithJSON(w, 201, response)
}

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	var user CreatedUser
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&user); err != nil {
		respondWithError(w, 400, "Something when wrong")
		return
	}
	ctx := r.Context()
	data, err := cfg.db.CheckUserByEmail(ctx, user.Email)
	if err != nil {
		respondWithError(w, 401, "Incorrect email")
		return
	}
	match, err := auth.CheckPasswordHash(user.Password, data.HashedPassword)
	if err != nil || match == false {
		respondWithError(w, 401, "Wrong password")
		return
	}
	JWTtoken, err := auth.MakeJWT(data.ID, cfg.jwt)
	if err != nil {
		respondWithError(w, 400, "No Token")
		return
	}
	refreshToken := auth.MakeRefreshToken()

	params := database.CreateRefreshTokenParams{
		Token:     refreshToken,
		UserID:    data.ID,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 60),
	}
	dataToken, err := cfg.db.CreateRefreshToken(ctx, params)
	if err != nil {
		respondWithError(w, 400, "could not make refresh token")
		return
	}
	response := FoundUser{
		ID:           data.ID,
		CreatedAt:    data.CreatedAt,
		UpdatedAt:    data.UpdatedAt,
		Email:        data.Email,
		Token:        JWTtoken,
		RefreshToken: dataToken.Token,
	}
	respondWithJSON(w, 200, response)
}

func (cfg *apiConfig) handlerUsersUpdate(w http.ResponseWriter, r *http.Request) {

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "No token")
		return
	}
	userID, err := auth.ValidateJWT(token, cfg.jwt)
	if err != nil {
		respondWithError(w, 401, "No User found")
		return
	}
	ctx := r.Context()

	var update UpdateUser
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&update); err != nil {
		respondWithError(w, 400, "Something when wrong")
		return
	}
	if update.Email == "" || update.Password == "" {
		respondWithError(w, 400, "No Email or Password Given")
		return
	}
	hashedPass, err := auth.HashPassword(update.Password)
	if err != nil {
		respondWithError(w, 400, "Could not update password")
	}
	params := database.UpdateUserParams{
		ID:             userID,
		Email:          update.Email,
		HashedPassword: hashedPass,
	}
	data, err := cfg.db.UpdateUser(ctx, params)
	if err != nil {
		respondWithError(w, 500, "User could not be updated")
		return
	}
	response := FoundUser{
		ID:        userID,
		CreatedAt: data.CreatedAt,
		UpdatedAt: data.UpdatedAt,
		Email:     data.Email,
	}
	respondWithJSON(w, 200, response)
}
