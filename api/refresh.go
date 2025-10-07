package api

import (
	"net/http"
	"time"

	"github.com/MarianGheorghiu/chirpy/internal/auth"
)

func (cfg *APIConfig) HandlerRefresh(w http.ResponseWriter, r *http.Request) {
	strToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid or expired token", nil)
		return
	}

	refreshToken, err := cfg.Queries.GetRefreshToken(r.Context(), strToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid or expired token", nil)
		return
	}
	if refreshToken.RevokedAt.Valid {
		respondWithError(w, http.StatusUnauthorized, "invalid or expired token", nil)
		return
	}

	if time.Now().UTC().After(refreshToken.ExpiresAt) {
		respondWithError(w, http.StatusUnauthorized, "invalid or expired token", nil)
		return
	}

	user, err := cfg.Queries.GetUserFromRefreshToken(r.Context(), strToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't get user for refresh token", err)
		return
	}

	accessToken, err := auth.MakeJWT(user.ID, cfg.TokenSecret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error generating token", err)
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"token": accessToken})

}
