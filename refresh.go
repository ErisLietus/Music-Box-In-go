package main

import (
	"net/http"

	"github.com/ErisLietus/go-http-client/internal/auth"
)

type CreateRefreshToken struct {
	Token string `json:"token"`
}

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	bearerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "Could not use token")
		return
	}

	user, err := cfg.db.GetUserFromRefreshToken(ctx, bearerToken)
	if err != nil {
		respondWithError(w, 401, "Token has been revoked")
		return
	}

	newAccessToken, err := auth.MakeJWT(user.ID, cfg.jwt)
	if err != nil {
		respondWithError(w, 500, "Invalid")
		return
	}

	param := CreateRefreshToken{
		Token: newAccessToken,
	}
	respondWithJSON(w, 200, param)
}

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	bearerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "Could not use token")
		return
	}
	cfg.db.RevokeToken(ctx, bearerToken)

	respondWithJSON(w, 204, nil)

}
