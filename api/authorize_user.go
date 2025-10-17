package api

import (
	"encoding/json"
	"net/http"

	"github.com/MarianGheorghiu/chirpy/internal/auth"
	"github.com/MarianGheorghiu/chirpy/internal/database"
)

func (cfg *APIConfig) HandlerAuthorization(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid or missing token", nil)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.TokenSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid or missing token", nil)
		return
	}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	params := userParameters{}

	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	hashedPass, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't hash password", err)
		return
	}

	userParams := database.UpdateUserParams{
		ID:             userID,
		Email:          params.Email,
		HashedPassword: hashedPass,
	}

	dbUserUpdated, err := cfg.Queries.UpdateUser(r.Context(), userParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't update user", err)

		return
	}

	respondWithJSON(w, http.StatusOK, User{
		ID:          dbUserUpdated.ID,
		CreatedAt:   dbUserUpdated.CreatedAt,
		UpdatedAt:   dbUserUpdated.UpdatedAt,
		Email:       dbUserUpdated.Email,
		IsChirpyRed: dbUserUpdated.IsChirpyRed,
	})
}
