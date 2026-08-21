package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ErisLietus/Music_box_go/internal/auth"
	"github.com/ErisLietus/Music_box_go/internal/database"
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
		return
	}
	emailHash := auth.HashEmail(user.Email, os.Getenv("EMAILSECRET"))

	params := database.CreateUserParams{
		HashedEmail:    emailHash,
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
		Email:     data.HashedEmail,
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
	lookupHash := auth.HashEmail(user.Email, os.Getenv("EMAILSECRET"))
	data, err := cfg.db.CheckUserByEmail(ctx, lookupHash)
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
	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, 400, "could not make token")
		return
	}

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
		Email:        data.HashedEmail,
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
		return
	}
	params := database.UpdateUserParams{
		ID:             userID,
		HashedEmail:    update.Email,
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
		Email:     data.HashedEmail,
	}
	respondWithJSON(w, 200, response)
}
