package api

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *APIConfig) HandlerChirpGet(w http.ResponseWriter, r *http.Request) {
	stringId := r.PathValue("chirpID")
	id, err := uuid.Parse(stringId)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid chirp id", err)
		return
	}

	dbChrip, err := cfg.Queries.GetChirp(r.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "chirp not found", nil)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "could not fetch chirp", err)
		return
	}

	respondWithJSON(w, http.StatusOK, Chirp{
		ID:        dbChrip.ID,
		CreatedAt: dbChrip.CreatedAt,
		UpdatedAt: dbChrip.UpdatedAt,
		UserID:    dbChrip.UserID,
		Body:      dbChrip.Body,
	})

}
