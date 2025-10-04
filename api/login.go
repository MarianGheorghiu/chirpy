package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/MarianGheorghiu/chirpy/internal/auth"
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
		respondWithError(w, http.StatusUnauthorized, "email and password required", nil)
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

	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		},
	})
}
