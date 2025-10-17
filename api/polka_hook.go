package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/MarianGheorghiu/chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *APIConfig) HandlerPolkaWebhooks(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil || apiKey != cfg.PolkaKey {
		respondWithError(w, http.StatusUnauthorized, "invalid or missing api key", nil)
		return
	}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	params := polkaWebhook{}

	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid webhook payload", err)
		return
	}

	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	uid, err := uuid.Parse(strings.TrimSpace(params.Data.UserID))
	if err != nil {
		respondWithError(w, http.StatusNotFound, "user not found", nil)
		return
	}

	_, err = cfg.Queries.UpgradeUserToChirpyRed(r.Context(), uid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "user not found", nil)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "couldn't upgrade user", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
