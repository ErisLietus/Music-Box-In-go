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
	"github.com/google/uuid"
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

func (cfg *apiConfig) handlerDeletePlaylist(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
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
	ID := r.URL.Query().Get("playlist_id")

	playlistId, err := uuid.Parse(ID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid", err)
		return
	}
	playlist, err := cfg.db.GetPlaylistByID(ctx, playlistId)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid", err)
		return
	}

	if userID != playlist.UserID {
		respondWithError(w, http.StatusForbidden, "Invalid", fmt.Errorf("Forbidden"))
		return
	}
	if err := cfg.db.DeletePlaylist(ctx, playlistId); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not delete playlist", err)
		return
	}
	respondWithJSON(w, http.StatusOK, "Deleted")
}

func (cfg *apiConfig) handlerUpdatePlaylists(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
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
	ID := r.URL.Query().Get("playlist_id")

	playlistId, err := uuid.Parse(ID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid", err)
		return
	}
	playlist, err := cfg.db.GetPlaylistByID(ctx, playlistId)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid", err)
		return
	}

	if userID != playlist.UserID {
		respondWithError(w, http.StatusForbidden, "You can not update this", fmt.Errorf("Forbidden"))
		return
	}

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

	params := database.UpdatePlaylistParams{
		ID:               playlistId,
		Name:             req.Name,
		IsPublic:         req.IsPublic,
		AllowCollabEdits: req.AllowEdit,
	}

	newPlaylist, err := cfg.db.UpdatePlaylist(ctx, params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "could not update", err)
		return
	}
	respondWithJSON(w, http.StatusOK, newPlaylist)
}
