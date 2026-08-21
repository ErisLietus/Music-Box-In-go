package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/ErisLietus/Music_box_go/internal/auth"
	"github.com/ErisLietus/Music_box_go/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func setupTestDB(t *testing.T) *apiConfig {
	_ = godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		t.Skip("Skipping integration test: DB_URL not set")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		t.Fatalf("Failed to connect to test db: %v", err)
	}

	return &apiConfig{
		db:  database.New(db),
		jwt: "test_secret_key_12345678901234567",
	}
}

func TestCreatePlaylistFullIntegration(t *testing.T) {
	cfg := setupTestDB(t)
	ctx := t.Context()

	_ = cfg.db.DeleteUsers(ctx)

	t.Cleanup(func() {
		_ = cfg.db.DeleteUsers(ctx)
	})

	hashedPass, _ := auth.HashPassword("password123")
	user, err := cfg.db.CreateUser(ctx, database.CreateUserParams{
		Username:       "testuser_playlist",
		HashedEmail:    "test_playlist@example.com",
		HashedPassword: hashedPass,
	})
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	token, err := auth.MakeJWT(user.ID, cfg.jwt)
	if err != nil {
		t.Fatalf("Failed to make JWT: %v", err)
	}

	payload := []byte(`{"name": "Lo-Fi Beats", "is_public": true, "allow_collab_edits": false}`)
	req := httptest.NewRequest(http.MethodPost, "/api/playlists", bytes.NewBuffer(payload))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	cfg.handlerCreatePlaylist(rr, req)

	if rr.Code != http.StatusCreated && rr.Code != http.StatusOK {
		t.Fatalf("Expected 201/200, got %d. Body: %s", rr.Code, rr.Body.String())
	}

	var created database.Playlist
	if err := json.NewDecoder(rr.Body).Decode(&created); err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	if created.Name != "Lo-Fi Beats" {
		t.Errorf("Expected playlist name 'Lo-Fi Beats', got '%s'", created.Name)
	}
	if created.UserID != user.ID {
		t.Errorf("Expected playlist user ID to match %v, got %v", user.ID, created.UserID)
	}

	mediaPayload := []byte(`{
		"playlist_name": "Lo-Fi Beats",
		"title": "Midnight Chill Beat",
		"media_url": "https://www.youtube.com/watch?v=jfKfPfyJRdk"
	}`)

	mediaReq := httptest.NewRequest(http.MethodPost, "/api/media", bytes.NewBuffer(mediaPayload))
	mediaReq.Header.Set("Authorization", "Bearer "+token)
	mediaReq.Header.Set("Content-Type", "application/json")
	mediaRR := httptest.NewRecorder()

	cfg.handlerImportMediaLink(mediaRR, mediaReq)

	if mediaRR.Code != http.StatusCreated && mediaRR.Code != http.StatusOK {
		t.Fatalf("Expected 201/200 for media import, got %d. Body: %s", mediaRR.Code, mediaRR.Body.String())
	}

	mediaList, err := cfg.db.GetMediaByPlaylist(ctx, created.ID)
	if err != nil {
		t.Fatalf("Failed to fetch media from DB: %v", err)
	}

	if len(mediaList) != 1 {
		t.Fatalf("Expected 1 media item in playlist, got %d", len(mediaList))
	}

	if mediaList[0].Title != "Midnight Chill Beat" {
		t.Errorf("Expected title 'Midnight Chill Beat', got '%s'", mediaList[0].Title)
	}

	if mediaList[0].Position != 1 {
		t.Errorf("Expected position 1, got %d", mediaList[0].Position)
	}
}
