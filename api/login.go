package api

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/MarianGheorghiu/chirpy/internal/auth"
	"github.com/MarianGheorghiu/chirpy/internal/database"
)

func (cfg *APIConfig) HandlerLogin(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	params := userParameters{}
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't decode parameters", err)
		return
	}

	email := strings.TrimSpace(params.Email)
	password := strings.TrimSpace(params.Password)
	if email == "" || password == "" {
		respondWithError(w, http.StatusBadRequest, "email and password required", nil)
		return
	}

	user, err := cfg.Queries.GetUserByEmail(r.Context(), email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	ok, err := auth.CheckPasswordHash(password, user.HashedPassword)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	token, err := auth.MakeJWT(user.ID, cfg.TokenSecret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Token error", err)
		return
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Refresh Token error", err)
		return
	}

	now := time.Now().UTC()
	rtTTL := 60 * 24 * time.Hour
	expiresAt := now.Add(rtTTL)

	err = cfg.Queries.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:     refreshToken,
		UserID:    user.ID,
		ExpiresAt: expiresAt,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not persist refresh token", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:          user.ID,
			CreatedAt:   user.CreatedAt,
			UpdatedAt:   user.UpdatedAt,
			Email:       user.Email,
			IsChirpyRed: user.IsChirpyRed,
		},
		Token:        token,
		RefreshToken: refreshToken,
	})
}
