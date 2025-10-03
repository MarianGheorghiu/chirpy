package api

import "net/http"

func (cfg *APIConfig) HandlerChirpsList(w http.ResponseWriter, r *http.Request) {
	data, err := cfg.Queries.ListChirps(r.Context())
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

	respondWithJSON(w, http.StatusOK, chirps)
}
