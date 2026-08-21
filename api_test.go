package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandlerUsersCreateInvalidBody(t *testing.T) {
	// 1. Create a dummy API config (no DB needed for syntax validation test)
	cfg := &apiConfig{}

	// 2. Create invalid JSON payload
	invalidJSON := []byte(`{"invalid_json": `)

	req := httptest.NewRequest(http.MethodPost, "/api/users", bytes.NewBuffer(invalidJSON))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	// 3. Call the handler
	cfg.handlerUsersCreate(rr, req)

	// 4. Assert 400 Bad Request
	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}

	// 5. Inspect the JSON error response
	var res errorResponse
	err := json.NewDecoder(rr.Body).Decode(&res)
	if err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if res.Error == "" {
		t.Errorf("expected an error message in response body")
	}
}
