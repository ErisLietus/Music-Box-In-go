package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
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
	fmt.Printf("playlist endpoint was hit")
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized", err)
		return
	}
	userID, err := auth.ValidateJWT(token, cfg.jwt)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid token", err)
		return
	}

	// 2. Decode the request body
	var req CreatePlaylistRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	if req.Name == "" {
		respondWithError(w, http.StatusBadRequest, "Playlist name is required", err)
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
		respondWithError(w, http.StatusConflict, "A playlist with this name already exists", err)
		return
	}
	if err != nil {
		log.Printf("CreatePlaylist DB Error: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Could not create playlist", err)
		return
	}

	// 4. Return the created playlist object
	respondWithJSON(w, http.StatusCreated, newPlaylist)
}

func (cfg *apiConfig) handlerGetPlaylists(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("get playlists was hit")
	ctx := r.Context()
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid", err)
		return
	}
	userID, err := auth.ValidateJWT(token, cfg.jwt)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid", err)
		return
	}
	playlists, err := cfg.db.GetUserplaylists(ctx, userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't get playlists", err)
	}

	respondWithJSON(w, http.StatusOK, playlists)
}
