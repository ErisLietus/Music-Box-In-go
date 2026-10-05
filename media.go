package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/ErisLietus/Music_box_go/internal/auth"
	"github.com/ErisLietus/Music_box_go/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerImportMediaLink(w http.ResponseWriter, r *http.Request) {
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

	var req ImportMediaRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload", err)
		return
	}
	ID := r.URL.Query().Get("playlist_id")

	playlistId, err := uuid.Parse(ID)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid", err)
	}

	if req.Title == "" || req.MediaURL == "" {
		respondWithError(w, http.StatusBadRequest, "Missing information to complete action please try again", err)
		return
	}

	playlist, err := cfg.db.GetPlaylistByID(ctx, playlistId)

	if playlist.UserID != userID || playlist.AllowCollabEdits == false {
		respondWithError(w, http.StatusBadRequest, "Edit not allowed", err)
	}

	if err != nil {
		respondWithError(w, http.StatusNotFound, "Playlist not found", err)
		return
	}
	maxPosition, err := cfg.db.GetMaxMediaPosition(ctx, playlist.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not get media position", err)
		return
	}
	nextPosition := maxPosition + 1

	mediaType, err := LinkMediaCheck(req.MediaURL)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Only youtube links are supported", err)
		return
	}

	embedUrl, err := adjustMediaURL(req.MediaURL)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "this is not a valid url", err)
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
		respondWithError(w, http.StatusInternalServerError, "Could not create media", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, createdMedia)

}

func LinkMediaCheck(rawURL string) (database.MediaType, error) {
	parsedURL, err := url.ParseRequestURI(rawURL)
	if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
		return "", fmt.Errorf("invalid URL: must include scheme and host")
	}

	host := strings.ToLower(parsedURL.Host)
	if strings.Contains(host, "youtube.com") || strings.Contains(host, "youtu.be") {
		return database.MediaTypeLink, nil
	}
	return "", fmt.Errorf("only youtube links are supported")
}

func adjustMediaURL(URL string) (string, error) {
	parsedURL, err := url.ParseRequestURI(URL)
	if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
		return "", fmt.Errorf("invalid URL: must include scheme and host")
	}
	videoID := parsedURL.Query().Get("v")
	finalURL := fmt.Sprintf("https://www.youtube.com/embed/%s", videoID)

	return finalURL, nil
}

type ImportMediaRequest struct {
	MediaURL     string `json:"media_url"`
	PlaylistName string `json:"playlist_name"`
	Title        string `json:"title"`
}

func (cfg *apiConfig) uploadMediaMP3(w http.ResponseWriter, r *http.Request) {
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

	fmt.Println("Upload file handler hit")

	if err := r.ParseMultipartForm(100 << 20); err != nil {
		respondWithError(w, http.StatusBadRequest, "Could not parse form", err)
		return
	}

	playlistName := r.FormValue("playlist_name")
	title := r.FormValue("title")

	file, handler, err := r.FormFile("uploadedFile")
	if err != nil {
		fmt.Println("Could not retrieve file")
		fmt.Println(err)
		return
	}
	defer file.Close()
	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	tempFile, err := os.CreateTemp("media-storage", "upload-*.mp3")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not upload", err)
		return
	}
	defer tempFile.Close()

	fmt.Printf("filename: %+v\n", tempFile.Name())

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not create media", err)
		return
	}
	_, err = tempFile.Write(fileBytes)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not create media", err)
		return
	}
	fmt.Println("Successful upload")

	mediaType, err := detectMediaTypeFromFilename(handler.Filename)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "File  is invalid", err)
		return
	}
	if mediaType != database.MediaTypeAudio && mediaType != database.MediaTypeVideo {
		respondWithError(w, http.StatusBadRequest, "Format is not a valid audio/video file", err)
		return
	}

	playlist, err := cfg.db.GetPlaylistByUser(ctx, database.GetPlaylistByUserParams{
		UserID: userID,
		Name:   playlistName,
	})
	if err != nil {
		respondWithError(w, http.StatusNotFound, "You cannot edit playlist that is not yours", err)
		return
	}
	maxPosition, err := cfg.db.GetMaxMediaPosition(ctx, playlist.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not get media position", err)
		return
	}
	nextPosition := maxPosition + 1

	param := database.CreateMediaParams{
		PlaylistID:    playlist.ID,
		Title:         title,
		FileUrl:       tempFile.Name(),
		Type:          mediaType,
		Position:      nextPosition,
		AddedByUserID: userID,
	}

	cfg.db.NoneLinkAdded(ctx, playlist.ID)

	createdMedia, err := cfg.db.CreateMedia(ctx, param)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not create media", err)
		return
	}
	respondWithJSON(w, 200, createdMedia)
}

func detectMediaTypeFromFilename(filename string) (database.MediaType, error) {
	lower := strings.ToLower(filename)
	if strings.HasSuffix(lower, ".mp3") || strings.HasSuffix(lower, ".ogg") || strings.HasSuffix(lower, ".wav") {
		return database.MediaTypeAudio, nil
	}
	if strings.HasSuffix(lower, ".mp4") || strings.HasSuffix(lower, ".webm") {
		return database.MediaTypeVideo, nil
	}
	return "", fmt.Errorf("unsupported file extension")
}

func (cfg *apiConfig) handlerGetMediaByPlaylist(w http.ResponseWriter, r *http.Request) {
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
	}

	media, err := cfg.db.GetMediaByPlaylist(ctx, playlistId)
	if err != nil {
		respondWithError(w, 404, "Not valid playlist", err)
		return
	}
}
