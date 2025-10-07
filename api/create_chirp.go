package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"unicode/utf8"

	"github.com/MarianGheorghiu/chirpy/internal/auth"
	"github.com/MarianGheorghiu/chirpy/internal/database"
)

func (cfg *APIConfig) HandlerChirpCreate(w http.ResponseWriter, r *http.Request) {
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
	params := chirpCreateInput{}

	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	cleaned, err := validateChirp(params.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error(), err)
		return
	}

	dbChirp, err := cfg.Queries.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   cleaned,
		UserID: userID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create chrip", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, Chirp{
		ID:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body:      dbChirp.Body,
		UserID:    dbChirp.UserID,
	})
}

func validateChirp(body string) (string, error) {
	const maxChirpLength = 140

	body = strings.TrimSpace(body)
	if body == "" {
		return "", errors.New("Chirp cannot be empty")
	}
	if utf8.RuneCountInString(body) > maxChirpLength {
		return "", errors.New("Chirp is too long")
	}

	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}
	cleaned := getCleanedBody(body, badWords)
	return cleaned, nil
}

func getCleanedBody(body string, badWords map[string]struct{}) string {
	words := strings.Split(body, " ")
	for i, word := range words {
		loweredWord := strings.ToLower(word)
		if _, ok := badWords[loweredWord]; ok {
			words[i] = "****"
		}
	}
	cleaned := strings.Join(words, " ")
	return cleaned
}
