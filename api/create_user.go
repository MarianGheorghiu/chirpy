package api

import (
	"encoding/json"

	"net/http"

	"github.com/MarianGheorghiu/chirpy/internal/auth"
	"github.com/MarianGheorghiu/chirpy/internal/database"
)

func (cfg *APIConfig) HandlerUsersCreate(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	decoder := json.NewDecoder(r.Body)
	params := userParameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	hashed, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not hash password", err)
		return
	}

	newUser := database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: hashed,
	}

	user, err := cfg.Queries.CreateUser(r.Context(), newUser)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, response{
		User: User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		},
	})

}
