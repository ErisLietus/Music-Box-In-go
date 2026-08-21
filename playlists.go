package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/ErisLietus/Music_box_go/internal/auth"
	"github.com/ErisLietus/Music_box_go/internal/database"
)

type CreatePlaylistRequest struct {
	Name      string `json:"name"`
	IsPublic  bool   `json:"is_public"`
	AllowEdit bool   `json:"allow_collab_edits"`
}

func (cfg *apiConfig) handlerCreatePlaylist(w http.ResponseWriter, r *http.Request) {
	// 1. Authenticate user from Bearer token
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	userID, err := auth.ValidateJWT(token, cfg.jwt)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid token")
		return
	}

	// 2. Decode the request body
	var req CreatePlaylistRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if req.Name == "" {
		respondWithError(w, http.StatusBadRequest, "Playlist name is required")
		return
	}

	// 3. Insert into Database
	params := database.CreatePlaylistParams{
		UserID:           userID,
		Name:             req.Name,
		IsPublic:         req.IsPublic,
		AllowCollabEdits: req.AllowEdit,
	}

	ctx := r.Context()
	newPlaylist, err := cfg.db.CreatePlaylist(ctx, params)
	if errors.Is(err, sql.ErrNoRows) {
		respondWithError(w, http.StatusConflict, "A playlist with this name already exists")
		return
	}
	if err != nil {
		log.Printf("CreatePlaylist DB Error: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Could not create playlist")
		return
	}

	// 4. Return the created playlist object
	respondWithJSON(w, http.StatusCreated, newPlaylist)
}
