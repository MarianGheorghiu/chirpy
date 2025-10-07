package api

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/MarianGheorghiu/chirpy/internal/auth"
	"github.com/MarianGheorghiu/chirpy/internal/database"
)

func (cfg *APIConfig) HandlerRevoke(w http.ResponseWriter, r *http.Request) {
	strToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid or expired token", nil)
		return
	}

	err = cfg.Queries.RevokeRefreshToken(r.Context(), database.RevokeRefreshTokenParams{
		Token: strToken,
		RevokedAt: sql.NullTime{
			Time:  time.Now().UTC(),
			Valid: true,
		},
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not revoke token", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
