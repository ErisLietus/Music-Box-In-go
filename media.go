package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/ErisLietus/Music_box_go/internal/auth"
	"github.com/ErisLietus/Music_box_go/internal/database"
)

type ImportMediaRequest struct {
	MediaURL     string `json:"media_url"`
	PlaylistName string `json:"playlist_name"`
	Title        string `json:"title"`
}

func (cfg *apiConfig) handlerImportMediaLink(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
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
	var req ImportMediaRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if req.Title == "" || req.PlaylistName == "" || req.MediaURL == "" {
		respondWithError(w, http.StatusBadRequest, "Missing information to complete action please try again")
		return
	}

	playlist, err := cfg.db.GetPlaylistByUser(ctx, database.GetPlaylistByUserParams{
		UserID: userID,
		Name:   req.PlaylistName,
	})
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Playlist not found")
		return
	}
	maxPosition, err := cfg.db.GetMaxMediaPosition(ctx, playlist.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not get media position")
		return
	}
	nextPosition := maxPosition + 1

	mediaType, err := detectMediaType(req.MediaURL)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Url is invalid")
		return
	}
	if mediaType != database.MediaTypeLink {
		respondWithError(w, http.StatusBadRequest, "Format is not a link ")
		return
	}

	embedUrl, err := adjustMediaURL(req.MediaURL)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "this is not a valid url")
		return
	}

	param := database.CreateMediaParams{
		PlaylistID:    playlist.ID,
		Title:         req.Title,
		FileUrl:       embedUrl,
		Type:          mediaType,
		Position:      nextPosition,
		AddedByUserID: userID,
	}
	createdMedia, err := cfg.db.CreateMedia(ctx, param)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not create media")
		return
	}

	respondWithJSON(w, http.StatusCreated, createdMedia)

}

func detectMediaType(rawURL string) (database.MediaType, error) {
	parsedURL, err := url.ParseRequestURI(rawURL)
	if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
		return "", fmt.Errorf("invalid URL: must include scheme and host")
	}

	host := strings.ToLower(parsedURL.Host)
	if strings.Contains(host, "youtube.com") || strings.Contains(host, "youtu.be") {
		return database.MediaTypeLink, nil
	}

	path := strings.ToLower(parsedURL.Path)
	if strings.HasSuffix(path, ".mp3") || strings.HasSuffix(path, ".ogg") || strings.HasSuffix(path, ".wav") {
		return database.MediaTypeAudio, nil
	}
	if strings.HasSuffix(path, ".mp4") || strings.HasSuffix(path, ".webm") {
		return database.MediaTypeVideo, nil
	}
	return database.MediaTypeLink, nil
}

func adjustMediaURL(URL string) (string, error) {
	parsedURL, err := url.ParseRequestURI(URL)
	if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
		return "", fmt.Errorf("invalid URL: must include scheme and host")
	}
	parsedURL.Path = "embed/"
	parsedURL.RawQuery = strings.TrimSpace(strings.TrimPrefix(parsedURL.RawQuery, "v="))
	splitURL := strings.Split(parsedURL.String(), "?")
	if len(splitURL) == 3 {
		splitURL[1] = ""
	}
	finalURL := strings.Join(splitURL, "")
	fmt.Println(finalURL)

	return finalURL, nil
}
