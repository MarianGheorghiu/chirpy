package api

import (
	"net/http"
	"sort"
	"strings"

	"github.com/MarianGheorghiu/chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *APIConfig) HandlerChirpsList(w http.ResponseWriter, r *http.Request) {
	authorIDStr := r.URL.Query().Get("author_id")
	sortParam := strings.ToLower(r.URL.Query().Get("sort"))

	var (
		data []database.Chirp
		err  error
	)

	if authorIDStr == "" {
		data, err = cfg.Queries.ListChirps(r.Context())
	} else {
		authorID, parseErr := uuid.Parse(authorIDStr)
		if parseErr != nil {
			respondWithError(w, http.StatusBadRequest, "invalid author_id", parseErr)
			return
		}
		data, err = cfg.Queries.ListChirpsByAuthorID(r.Context(), authorID)
	}

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't receive chirps", err)
		return
	}

	chirps := make([]Chirp, 0, len(data))
	for _, c := range data {
		chirps = append(chirps, Chirp{
			ID:        c.ID,
			CreatedAt: c.CreatedAt,
			UpdatedAt: c.UpdatedAt,
			UserID:    c.UserID,
			Body:      c.Body,
		})
	}

	if sortParam == "desc" {
		sort.Slice(chirps, func(i, j int) bool {
			return chirps[i].CreatedAt.After(chirps[j].CreatedAt)
		})
	} else {
		sort.Slice(chirps, func(i, j int) bool {
			return chirps[i].CreatedAt.Before(chirps[j].CreatedAt)
		})
	}

	respondWithJSON(w, http.StatusOK, chirps)
}
