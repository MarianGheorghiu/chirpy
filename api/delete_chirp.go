package api

import (
	"database/sql"
	"net/http"

	"github.com/MarianGheorghiu/chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *APIConfig) HandlerChirpDelete(w http.ResponseWriter, r *http.Request) {
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

	chirpIDStr := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(chirpIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid chirp id", err)
		return
	}

	dbChirp, err := cfg.Queries.GetChirp(r.Context(), chirpID)
	if err != nil {
		if err == sql.ErrNoRows {
			respondWithError(w, http.StatusNotFound, "chirp not found", nil)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "could not fetch chirp", err)
		return
	}

	if userID != dbChirp.UserID {
		respondWithError(w, http.StatusForbidden, "not allowed", nil)
		return
	}

	if err := cfg.Queries.DeleteChirp(r.Context(), chirpID); err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't delete chirp", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
