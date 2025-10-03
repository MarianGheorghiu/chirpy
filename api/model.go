package api

import (
	"sync/atomic"

	"github.com/MarianGheorghiu/chirpy/internal/database"
)

type APIConfig struct {
	fileserverHits atomic.Int32
	Queries        *database.Queries
}

type chirpParams struct {
	Body string `json:"body"`
}
type errorResp struct {
	Error string `json:"error"`
}

type cleanedResp struct {
	CleanedBody string `json:"cleaned_body"`
}
