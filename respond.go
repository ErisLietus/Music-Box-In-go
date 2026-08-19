package main

import (
	"encoding/json"
	"net/http"
)

type errorResponse struct {
	Error string `json:"error"`
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	payload := errorResponse{
		Error: msg,
	}
	respondWithJSON(w, code, payload)

}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "Application/json; charset=utf-8")
	jsonData, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Error"))
		return
	}
	w.WriteHeader(code)
	w.Write(jsonData)

}
